package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user types.User) (int, error)
	CreateAdmin(user types.User) (int, bool, error)
	GetUser(username, password string) (types.User, error)
}

type Banners interface {
	CreateBanner(banner types.BannerRequest, tags []int) (int, error)
	GetBanner(featureID, tagID, limit, offset int) ([]types.BannerResponse, error)
	GetUserBanner(featureID, tagID int) (types.Content, error)
	DeleteBanner(bannerId int) error
	UpdateBanner(inputUpdate types.BannerRequest) error
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

type Repository struct {
	Authorization
	Banners
	Cache
}

func NewRepository(db *sqlx.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banners:       NewBannersPostgres(db),
		Cache:         NewRedisCache(redisClient)}
}
