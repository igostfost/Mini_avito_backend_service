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

func (h *Handler) GetBannersHandler(c *gin.Context) {

	isAdmin, ok := c.Get("isAdmin")
	if !isAdmin.(bool) || !ok {
		NewErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа - Только администраторы могут просматривать все баннеры")
		return
	}

	var input types.GetInputBanners
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	banners, err := h.utils.GetBanner(input.FeatureId, input.TagIds, input.Limit, input.Offset)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if banners == nil {
		NewErrorResponse(c, http.StatusNotFound, "Баннер не найден")
		return
	}

	c.JSON(http.StatusOK, banners)
}

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

	fmt.Println("получили данные напрямую из дб")
	c.JSON(http.StatusOK, userBanner)
}

func (h *Handler) GetUserBannerHandler(c *gin.Context) {
	// Привязываем параметры из URL запроса
	var input types.GetInputBanners
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}
	if input.UseLastRevision {
		h.GetUserBannerFromDB(c, input.TagIds, input.FeatureId)
		// fmt.Println("Ушли получать напрямую из бд")
		return
	}

	// Формируем ключ для кеша Redis
	cacheKey := fmt.Sprintf("user_banner:%d:%d", input.FeatureId, input.TagIds)

	// Проверяем, есть ли данные в кеше Redis
	cachedData, err := h.utils.Get(c, cacheKey)
	if err == nil {
		// Данные найдены в кеше, возвращаем их пользователю
		var userBanner types.Content
		cachedDataString, ok := cachedData.(string)
		if !ok {
			// Обработка ошибки, если cachedData не является строкой
			NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении данных из кеша")
			return
		}
		if err := json.Unmarshal([]byte(cachedDataString), &userBanner); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при разборе данных из кеша")
			return
		}
		// fmt.Println("Считали из кеша")
		c.JSON(http.StatusOK, userBanner)
		return
	}

	// Данных нет в кеше, получаем их из базы данных
	userBanner, err := h.utils.GetUserBanner(input.FeatureId, input.TagIds)
	// fmt.Println("Данных нет в кеше, получаем их из базы данных")
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			NewErrorResponse(c, http.StatusNotFound, "Баннер не найден 1")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	// Сохраняем данные в кеш Redis
	jsonData, err := json.Marshal(userBanner)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при сериализации данных")
		return
	}

	if err := h.utils.Set(c, cacheKey, string(jsonData)); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении данных в кеш")
		return
	}
	// fmt.Println("Сохранили данные в кеш")

	// Возвращаем данные пользователю
	c.JSON(http.StatusOK, userBanner)
}

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

	// c.JSON(http.StatusOK, inputUpdate)
}

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
