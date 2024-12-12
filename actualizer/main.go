package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test-work/internal/config"
	"test-work/internal/repositories/banner"
	"test-work/internal/repositories/banner_cached"
	"test-work/internal/services"
	banner_service "test-work/internal/services/banner"
	pgStore "test-work/internal/storages/postgres"
	redisStore "test-work/internal/storages/redis"
	"time"
)

func main() {
	os.Setenv("TZ", "UTC")
	ctx := context.Background()
	cfg := config.SLoad()

	rdb, err := redisStore.New(ctx, cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	pdb, err := pgStore.New(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	bannerRepo := banner.NewRepository(pdb)
	bannerCachedRepo := banner_cached.NewRepository(rdb)

	serviceLayer := &services.ServiceLayer{BannerService: &banner_service.Service{
		BannerRepo:       bannerRepo,
		BannerCachedRepo: bannerCachedRepo,
	}}

	fmt.Printf("Actualizer starting...\n")

	// Создаём новый планировщик задач
	c := cron.New()

	// Добавляем задачи
	// Задача будет выполняться каждую минуту
	_, err = c.AddFunc("* * * * *", func() {
		fmt.Println("Задача выполняется каждую минуту:", time.Now())
		serviceLayer.BannerService.ActualizeStats(ctx)
	})
	if err != nil {
		fmt.Println("Ошибка добавления задачи:", err)
	}

	c.Start()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Println("stopping Actualizer", slog.String("signal", sign.String()))
}
