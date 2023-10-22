package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(Addr, Password string, DB int) CacheStore {
	return &RedisCache{client: redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	}),
	}

}
func (cache *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := cache.client.Get(ctx, "key").Result()
	if err != nil {
		return result, fmt.Errorf("Failed to get key:", err)
	}
	return result, nil
}
// 增加键值对
func (cache *RedisCache) Add(ctx context.Context, key string, value interface{}, expires time.Duration) error {
	err := cache.client.Set(ctx, "key", "value", expires).Err()
	if err != nil {
		return fmt.Errorf("Failed to add key:", err)
	}
	return nil
}
func (cache *RedisCache) Delete(ctx context.Context, key string) error {
	err := cache.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("Failed to delete key:", err)
	}
	return nil
}

// 更改value值
func (cache *RedisCache) Update(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	err := cache.client.Set(ctx, "key", "new value", expire).Err()
	if err != nil {
		return fmt.Errorf("Failed to change value:", err)
	}
	return nil
}
