package cache

import "context"

type Repository interface {
	SetCachedData(ctx context.Context, key string, value interface{}) error
	GetCachedData(ctx context.Context, key string) string
}
