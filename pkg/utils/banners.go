package utils

import (
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
)

type BannersUtils struct {
	repo repository.Banners
}

func NewBannersUtils(repo repository.Banners) *BannersUtils {
	return &BannersUtils{repo: repo}
}

func (u *BannersUtils) CreateBanner(banner types.BannerRequest, tags []int) (int, error) {
	return u.repo.CreateBanner(banner, tags)
}

func (u *BannersUtils) GetBanner(featureID, tagID, limit, offset int) ([]types.BannerResponse, error) {
	return u.repo.GetBanner(featureID, tagID, limit, offset)
}

func (u *BannersUtils) GetUserBanner(featureID, tagID int) (types.Content, error) {
	return u.repo.GetUserBanner(featureID, tagID)
}

func (u *BannersUtils) DeleteBanner(bannerId int) error {
	return u.repo.DeleteBanner(bannerId)
}

func (u *BannersUtils) UpdateBanner(banner types.BannerRequest) error {
	return u.repo.UpdateBanner(banner)
}
