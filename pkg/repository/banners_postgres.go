package repository

import (
	"database/sql"
	"fmt"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"github.com/jmoiron/sqlx"
	"strings"
)

type BannersPostgres struct {
	db *sqlx.DB
}

func NewBannersPostgres(db *sqlx.DB) *BannersPostgres {
	return &BannersPostgres{db: db}
}

func (r *BannersPostgres) CreateBanner(banner types.BannerRequest, tags []int) (int, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Вставляем запись о баннере в таблицу banners
	bannerID, err := r.insertBanner(tx, banner)
	if err != nil {
		return 0, err
	}

	// Создаем связи между баннером и тегами в таблице banner_tags
	for _, tagID := range tags {

		err = r.insertBannerTag(tx, bannerID, tagID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, tx.Commit()
}

// Вставляет запись о баннере в таблицу banners
func (r *BannersPostgres) insertBanner(tx *sqlx.Tx, banner types.BannerRequest) (int, error) {

	query := fmt.Sprintf("INSERT INTO %s (feature_id, title, text, url, is_active)VALUES ($1, $2, $3, $4, $5) RETURNING banner_id", bannersTable)
	var bannerID int

	err := tx.QueryRow(query, banner.FeatureId, banner.Content.Title, banner.Content.Text, banner.Content.Url, banner.IsActive).Scan(&bannerID)
	if err != nil {
		return 0, err
	}
	return bannerID, nil
}

// Вставляет запись о связи баннера и тега в таблицу banner_tags
func (r *BannersPostgres) insertBannerTag(tx *sqlx.Tx, bannerID, tagID int) error {
	query := fmt.Sprintf("INSERT INTO %s (banner_id, tag_id)VALUES ($1, $2)", tagsTable)
	_, err := tx.Exec(query, bannerID, tagID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BannersPostgres) GetBanner(featureID, tagID, limit, offset int) ([]types.BannerResponse, error) {
	var banners []types.BannerResponse
	query := fmt.Sprintf("SELECT b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at FROM %s b", bannersTable)

	// Добавляем условия фильтрации по feature_id и tag_id
	var conditions []string
	if featureID != 0 {
		conditions = append(conditions, fmt.Sprintf("b.feature_id = %d", featureID))
	}
	if tagID != 0 {
		conditions = append(conditions, fmt.Sprintf("b.banner_id IN (SELECT banner_id FROM banner_tags WHERE tag_id = %d)", tagID))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Добавляем лимит и оффсет
	if limit != 0 || offset != 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	// Выполняем запрос
	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерируемся по результатам и формируем список баннеров
	for rows.Next() {
		var banner types.BannerResponse
		if err := rows.Scan(&banner.BannerId, &banner.FeatureId, &banner.Content.Title, &banner.Content.Text, &banner.Content.Url, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}
		banner.TagIds = tagID
		banners = append(banners, banner)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return banners, nil
}

func (r *BannersPostgres) GetUserBanner(featureID, tagID int) (types.Content, error) {
	var UserBanner types.Content

	query := fmt.Sprintf("SELECT b.title, b.text, b.url, b.is_active FROM %s b "+
		"JOIN %s bt ON bt.banner_id = b.banner_id WHERE bt.tag_id = %d AND b.feature_id = %d", bannersTable, tagsTable, tagID, featureID)

	// Выполняем запрос
	rows, err := r.db.Queryx(query)
	if err != nil {
		return UserBanner, err
	}
	defer rows.Close()

	if rows.Next() {
		var isActive bool
		if err := rows.Scan(&UserBanner.Title, &UserBanner.Text, &UserBanner.Url, &isActive); err != nil {
			return UserBanner, err
		}
		if !isActive {
			return UserBanner, sql.ErrNoRows
		}
	} else {
		return UserBanner, sql.ErrNoRows
	}

	return UserBanner, nil
}

func (r *BannersPostgres) DeleteBanner(bannerId int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE banner_id = $1", bannersTable)
	res, err := r.db.Exec(query, bannerId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
