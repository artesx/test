package banner

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"math/rand"
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
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(msg.Value))
		if err != nil {
			fmt.Println("error time parsed:", err)
		}

		key := fmt.Sprintf("%d:%s", parsedTime.Unix(), bannerID)
		err = h.BannerService.BannerCachedRepo.IncrToBanner(sess.Context(), key)
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

	formattedTime := roundedTime.Format("2006-01-02 15:04:05")

	msg := &sarama.ProducerMessage{
		Topic: "banner-clicks",
		Key:   sarama.StringEncoder(fmt.Sprintf("%s-%d", bannerID, rand.Intn(1000))),
		Value: sarama.StringEncoder(formattedTime),
	}

	_, _, err := s.Producer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}
