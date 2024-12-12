package banner

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"test-work/internal/repositories/banner"
	"test-work/internal/repositories/banner_cached"
	"time"
)

type Service struct {
	BannerRepo       *banner.Repository
	BannerCachedRepo *banner_cached.Repository
	Producer         sarama.SyncProducer
	ConsumerGroup    sarama.ConsumerGroup
}

type ConsumerGroupHandler struct {
	BannerService *Service
}

type GetStatisticBody struct {
	TsFrom string `json:"ts_from"`
	TsTo   string `json:"ts_to"`
}

type StatsResponse struct {
	BannerID uint64 `json:"banner_id"`
	Count    uint64 `json:"count"`
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer Group start")
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer Group end")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		bannerID := strings.Split(string(msg.Key), "-")[0]

		key := fmt.Sprintf("%s:%s", string(msg.Value), bannerID)

		err := h.BannerService.BannerCachedRepo.IncrToBanner(sess.Context(), key)
		fmt.Println("Incr to", bannerID)
		if err != nil {
			fmt.Println("error add to redis", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (s *Service) ConsumeBanner(ctx context.Context) {
	handler := &ConsumerGroupHandler{
		BannerService: s,
	}

	for {
		if err := s.ConsumerGroup.Consume(ctx, []string{"banner-clicks"}, handler); err != nil {
			log.Fatalf("Ошибка во время потребления: %v", err)
		}
	}
}

func (s *Service) ProduceBanner(bannerID string) (bool, error) {
	now := time.Now()
	roundedTime := now.Truncate(time.Minute)

	msg := &sarama.ProducerMessage{
		Topic: "banner-clicks",
		Key:   sarama.StringEncoder(fmt.Sprintf("%s-%d", bannerID, rand.Intn(1000))),
		Value: sarama.StringEncoder(fmt.Sprintf("%d", roundedTime.Unix())),
	}

	_, _, err := s.Producer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (s *Service) GetStatistic(bannerID uint64, body *GetStatisticBody) (*StatsResponse, error) {
	sum, err := s.BannerRepo.GetStatsSum(bannerID, body.TsFrom, body.TsTo)
	if err != nil {
		return nil, err
	}
	return &StatsResponse{
		Count:    sum,
		BannerID: bannerID,
	}, nil
}

func (s *Service) ActualizeStats(ctx context.Context) {
	result, err := s.BannerCachedRepo.GetAllCachedBanners(ctx)
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now()
	roundedTime := now.Truncate(time.Minute)

	for key, count := range result {
		arrKey := strings.Split(key, ":")
		timestamp, err := strconv.ParseInt(arrKey[0], 10, 64)
		if err != nil {
			log.Fatal("error parsed:", err)
		}
		bannerID, err := strconv.ParseInt(arrKey[1], 10, 64)
		if err != nil {
			log.Fatal("error parsed:", err)
		}
		formattedTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		if timestamp < roundedTime.Unix() {
			err = s.BannerRepo.AddStats(uint64(bannerID), formattedTime, count)
			if err != nil {
				fmt.Println("error add stats", err)
			}
			err = s.BannerCachedRepo.Delete(ctx, key)
			if err != nil {
				fmt.Println("error delete cache", err)
			}
		}
		fmt.Println("key:", key, "count:", count, "bannerID:", bannerID, "roundedTime:", roundedTime.Unix(), "timestamp:", formattedTime)
	}
}
