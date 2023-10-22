package cache

import (
	"context"
	"time"
)

type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)

	Add(ctx context.Context, key string, value interface{}, expires time.Duration) error
	Delete(ctx context.Context, key string) error
	Update(ctx context.Context, key string, value interface{}, expire time.Duration) error
}
