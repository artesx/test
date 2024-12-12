package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test-work/internal/config"
	"test-work/internal/http_transport"
	"test-work/internal/repositories/banner"
	"test-work/internal/services"
	banner_service "test-work/internal/services/banner"
	pgStore "test-work/internal/storages/postgres"
)

func main() {
	os.Setenv("TZ", "UTC")
	ctx := context.Background()
	cfg := config.SLoad()

	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Url}, conf)
	if err != nil {
		fmt.Println("kafka producer error:", err)
		log.Fatal(err)
	}
	defer producer.Close()

	pdb, err := pgStore.New(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	bannerRepo := banner.NewRepository(pdb)

	serviceLayer := &services.ServiceLayer{BannerService: &banner_service.Service{
		Producer:   producer,
		BannerRepo: bannerRepo,
	}}

	e := echo.New()

	api := http_transport.NewAPI(e)

	err = http_transport.RegisterHandlers(ctx, api, serviceLayer)

	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Println("stopping app", slog.String("signal", sign.String()))
}
