package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test-work/internal/config"
	"test-work/internal/repositories/banner_cached"
	"test-work/internal/services"
	banner_service "test-work/internal/services/banner"
	redisStore "test-work/internal/storages/redis"
)

func main() {
	ctx := context.Background()
	cfg := config.SLoad()

	// Конфигурация Consumer Group
	conf := sarama.NewConfig()
	conf.Version = sarama.V2_6_0_0 // Укажите версию Kafkax
	conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest // Читаем только новые сообщения

	// Создаем Consumer Group
	group, err := sarama.NewConsumerGroup([]string{cfg.Kafka.Url}, cfg.Kafka.GroupID, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	rdb, err := redisStore.New(ctx, cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	bannerCachedRepo := banner_cached.NewRepository(rdb)

	serviceLayer := &services.ServiceLayer{BannerService: &banner_service.Service{
		BannerCachedRepo: bannerCachedRepo,
		ConsumerGroup:    group,
	}}

	fmt.Printf("Consumer starting...\n")

	go serviceLayer.BannerService.ConsumeBanner(ctx)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Println("stopping consumer", slog.String("signal", sign.String()))
}
