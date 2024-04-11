package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"net/http"
	"strconv"
)

func (h *Handler) CreateBannerHandler(c *gin.Context) {
	//_, ok := c.Get("userId")
	//if !ok {
	//	NewErrorResponse(c, http.StatusInternalServerError, "user_id not found in context")
	//}
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

	c.JSON(http.StatusOK, banners)
}

func (h *Handler) GetUserBannerHandler(c *gin.Context) {

	var input types.GetInputBanners
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	userBanner, err := h.utils.GetUserBanner(input.FeatureId, input.TagIds)
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

	c.JSON(http.StatusOK, bannerId)
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
