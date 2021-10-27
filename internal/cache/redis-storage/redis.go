package redis_storage

import (
	"Yandex-Taxi-Clone/internal/cache"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const timeToLive = 10 * time.Minute

type redisRepository struct {
	storage *redis.Client
}

func New(storage *redis.Client) cache.Repository {
	return redisRepository{
		storage: storage,
	}
}

func (rr redisRepository) SetCachedData(ctx context.Context, key string, value interface{}) error {
	if err := rr.storage.Set(ctx, key, value, timeToLive).Err(); err != nil {
		return err
	}
	return nil
}

func (rr redisRepository) GetCachedData(ctx context.Context, key string) (string, error) {
	val, err := rr.storage.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
