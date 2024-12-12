package banner_cached

import (
	"context"
	"test-work/internal/storages/redis"
)

type Repository struct {
	Storage *redis.Storage
}

func NewRepository(storage *redis.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (r *Repository) IncrToBanner(ctx context.Context, key string) error {
	_, err := r.Storage.Db.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
