package utils

import (
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
)

type Authorization interface {
	CreateUser(user types.User) (int, error)
	CreateAdmin(user types.User) (int, bool, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, bool, error)
}

type Banners interface {
	CreateBanner(banner types.BannerRequest, tags []int) (int, error)
	GetBanner(featureID, tagID, limit, offset int) ([]types.BannerResponse, error)
	GetUserBanner(featureID, tagID int) (types.Content, error)
	DeleteBanner(bannerId int) error
	UpdateBanner(inputUpdate types.BannerRequest) error
}

type Utils struct {
	Authorization
	Banners
}

func NewUtils(repos *repository.Repository) *Utils {
	return &Utils{
		Authorization: NewAuthService(repos.Authorization),
		Banners:       NewBannersUtils(repos.Banners)}
}
