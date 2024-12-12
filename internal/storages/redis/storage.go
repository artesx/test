package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"test-work/internal/config"
)

type Storage struct {
	Db *redis.Client
}

func New(ctx context.Context, config config.RedisConfig) (*Storage, error) {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.Db,
	})

	err := rdb.Ping(ctx).Err()

	if err != nil {
		return nil, err
	}

	return &Storage{Db: rdb}, nil
}
