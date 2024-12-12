package banner_cached

import (
	"context"
	"strconv"
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

func (r *Repository) Delete(ctx context.Context, key string) error {
	err := r.Storage.Db.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllCachedBanners(ctx context.Context) (map[string]uint64, error) {
	var cursor uint64
	var keys []string
	result := make(map[string]uint64)
	for {
		var err error
		var newKeys []string
		newKeys, cursor, err = r.Storage.Db.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, newKeys...)

		if cursor == 0 {
			break
		}
	}

	values, err := r.Storage.Db.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		count, err := strconv.Atoi(values[i].(string))
		if err != nil {
			return nil, err
		}
		result[key] = uint64(count)
	}
	return result, nil
}
