package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
)

type RedisCache struct {
	repo repository.Cache
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{repo: repository.NewRedisCache(client)}
}

func (u *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return u.repo.Set(ctx, key, value)
}

func (u *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return u.repo.Get(ctx, key)
}

func (u *RedisCache) Delete(ctx context.Context, key string) error {
	return u.repo.Delete(ctx, key)
}
