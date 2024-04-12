package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"net/http"
)

// @Summary Регистрация
// @Tags auth
// @Description Создание учетной записи пользователя
// @ID create-user-account
// @Accept json
// @Produce json
// @Param input body types.User true "Информация об аккаунте"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input types.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.utils.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
}

// @Summary Регистрация администратора
// @Tags auth
// @Description Создание учетной записи администратора
// @ID create-admin-account
// @Accept json
// @Produce json
// @Param input body types.User true "Информация об аккаунте"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up/admin [post]
func (h *Handler) signUpAdmin(c *gin.Context) {
	var input types.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, isAdmin, err := h.utils.CreateAdmin(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      userId,
		"isAdmin": isAdmin,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Вход
// @Tags auth
// @Description Авторизация пользователя
// @ID login-user
// @Accept json
// @Produce json
// @Param input body signInInput true "Учетные данные"
// @Success 200 {string} string "Токен"
// @Failure 400,404 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.utils.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

// @Summary Вход администратора
// @Tags auth
// @Description Авторизация администратора
// @ID login-admin
// @Accept json
// @Produce json
// @Param input body signInInput true "Учетные данные"
// @Success 200 {string} string "Токен"
// @Failure 400,404 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in/admin [post]
func (h *Handler) signInAdmin(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.utils.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
