package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"net/http"
	"strconv"
)

// CreateBannerHandler создает новый баннер.
// @Summary Create Banner
// @Security ApiKeyAuth
// @Tags Banners
// @Description create Banner
// @ID create-banner
// @Accept  json
// @Produce  json
// @Param input body types.BannerRequest true "banner info"
// @Success 200 {integer} integer 1
// @Failure 400,403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /banner [post]
func (h *Handler) CreateBannerHandler(c *gin.Context) {

	isAdmin, ok := c.Get("isAdmin")
	if !isAdmin.(bool) || !ok {
		NewErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа - Только администраторы могут создавать баннеры")
		return
	}
	var inputBanner types.BannerRequest
	if err := c.BindJSON(&inputBanner); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	bannerId, err := h.utils.CreateBanner(inputBanner, inputBanner.TagIds)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"bannerId": bannerId,
	})

}

// GetBannersHandler возвращает список баннеров.
// @Summary Get Banners
// @Security ApiKeyAuth
// @Tags Banners
// @Description Получение списка баннеров
// @ID get-banners
// @Accept json
// @Produce json
// @Param input query types.GetInputBanners true "Параметры запроса"
// @Success 200 {array} types.BannerResponse "Список баннеров"
// @Failure 400,403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /banner [get]
func (h *Handler) GetBannersHandler(c *gin.Context) {
	var err error // Объявление переменной err здесь

	isAdmin, ok := c.Get("isAdmin")
	if !isAdmin.(bool) || !ok {
		NewErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа - Только администраторы могут просматривать все баннеры")
		return
	}

	featureIdStr := c.Query("feature_id")
	featureId := 0
	if featureIdStr != "" {
		featureId, err = strconv.Atoi(featureIdStr) // Присваивание значений err здесь
		if h.handleConversionError(c, err) {
			return
		}
	}

	tagIdStr := c.Query("tag_id")
	tagId := 0
	if tagIdStr != "" {
		tagId, err = strconv.Atoi(tagIdStr) // Присваивание значений err здесь
		if h.handleConversionError(c, err) {
			return
		}
	}

	limitStr := c.Query("limit")
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr) // Присваивание значений err здесь
		if h.handleConversionError(c, err) {
			return
		}
	}

	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr) // Присваивание значений err здесь
		if h.handleConversionError(c, err) {
			return
		}
	}

	banners, err := h.utils.GetBanner(featureId, tagId, limit, offset)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if banners == nil {
		NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		return
	}

	c.JSON(http.StatusOK, banners)
}

func (h *Handler) handleConversionError(c *gin.Context, err error) bool {
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return true
	}
	return false
}

// GetUserBannerFromDB возвращает баннер пользователя напрямую из бд с актуальной информацией
func (h *Handler) GetUserBannerFromDB(c *gin.Context, TagId, FeatureId int) {

	userBanner, err := h.utils.GetUserBanner(FeatureId, TagId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	c.JSON(http.StatusOK, userBanner)
}

// GetUserBannerHandler возвращает баннер пользователя.
// @Summary Get User Banner
// @Security ApiKeyAuth
// @Tags Banners
// @Description Получение баннера пользователя
// @ID get-user-banner
// @Accept json
// @Produce json
// @Param input query types.GetInputUserBanners true "Параметры запроса"
// @Success 200 {object} types.Content "Баннер пользователя"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user_banner [get]
func (h *Handler) GetUserBannerHandler(c *gin.Context) {

	tagId, err := strconv.Atoi(c.Query("tag_id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}
	featureId, err := strconv.Atoi(c.Query("feature_id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	useLastRevisionStr := c.Query("use_last_revision")
	useLastRevision := false
	if useLastRevisionStr != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if h.handleConversionError(c, err) {
			return
		}
	}

	//useLastRevision, err := strconv.ParseBool(c.Query("use_last_revision"))
	//if err != nil {
	//	NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
	//	return
	//}

	if useLastRevision {
		h.GetUserBannerFromDB(c, tagId, featureId)
		return
	}

	cacheKey := fmt.Sprintf("user_banner:%d:%d", featureId, tagId)

	cachedData, err := h.utils.Get(c, cacheKey)
	if err == nil {
		var userBanner types.Content
		cachedDataString, ok := cachedData.(string)
		if !ok {
			NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении данных из кеша")
			return
		}
		if err := json.Unmarshal([]byte(cachedDataString), &userBanner); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при разборе данных из кеша")
			return
		}
		c.JSON(http.StatusOK, userBanner)
		return
	}

	userBanner, err := h.utils.GetUserBanner(featureId, tagId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	jsonData, err := json.Marshal(userBanner)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при сериализации данных")
		return
	}

	if err := h.utils.Set(c, cacheKey, string(jsonData)); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении данных в кеш")
		return
	}

	c.JSON(http.StatusOK, userBanner)
}

// UpdateBannerHandler обновляет информацию о баннере.
// @Summary Update Banner
// @Security ApiKeyAuth
// @Tags Banners
// @Description Обновление информации о баннере
// @ID update-banner
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор баннера"
// @Param input body types.BannerRequest true "Информация о баннере"
// @Success 200 {string} string "OK"
// @Failure 400,403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /banner/{id} [patch]
func (h *Handler) UpdateBannerHandler(c *gin.Context) {
	isAdmin, ok := c.Get("isAdmin")
	if !isAdmin.(bool) || !ok {
		NewErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа - Только администраторы могут обновлять баннеры")
		return
	}

	bannerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	var inputUpdate types.BannerRequest
	if err := c.BindJSON(&inputUpdate); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}
	inputUpdate.BannerId = bannerId

	err = h.utils.UpdateBanner(inputUpdate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	c.JSON(http.StatusOK, "OK")

}

// DeleteBannerHandler удаляет баннер.
// @Summary Delete Banner
// @Security ApiKeyAuth
// @Tags Banners
// @Description Удаление баннера
// @ID delete-banner
// @Param id path int true "Идентификатор баннера"
// @Success 204 {string} string "Баннер успешно удален"
// @Failure 400,403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /banner/{id} [delete]
func (h *Handler) DeleteBannerHandler(c *gin.Context) {
	isAdmin, ok := c.Get("isAdmin")
	if !isAdmin.(bool) || !ok {
		NewErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа - Только администраторы могут удалять баннеры")
		return
	}
	bannerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
	}

	err = h.utils.DeleteBanner(bannerId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	c.JSON(http.StatusNoContent, "Баннер успешно удален")
}
