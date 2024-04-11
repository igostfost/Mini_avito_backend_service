package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/igostfost/avito_backend_trainee/pkg/utils"
)

type Handler struct {
	utils *utils.Utils
}

func NewHandler(utils *utils.Utils) *Handler {
	return &Handler{utils: utils}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in/admin", h.signInAdmin)
		auth.POST("/sign-up/admin", h.signUpAdmin)
	}

	banner := router.Group("/banner", h.userIdentity)
	{
		banner.POST("", h.CreateBannerHandler)       // Обработчик для создания баннера
		banner.GET("", h.GetBannersHandler)          // Обработчик для получения всех баннеров
		banner.PATCH("/:id", h.UpdateBannerHandler)  // Обработчик для обновления баннера
		banner.DELETE("/:id", h.DeleteBannerHandler) // Обработчик для удаления баннера
	}

	router.GET("/user_banner", h.userIdentity, h.GetUserBannerHandler) // Обработчик для получения баннера пользователя
	return router
}
