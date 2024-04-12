package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/igostfost/avito_backend_trainee/pkg/utils"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/igostfost/avito_backend_trainee/docs"
)

type Handler struct {
	utils *utils.Utils
}

func NewHandler(utils *utils.Utils) *Handler {
	return &Handler{utils: utils}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in/admin", h.signInAdmin)
		auth.POST("/sign-up/admin", h.signUpAdmin)
	}

	banner := router.Group("/banner", h.userIdentity)
	{
		banner.POST("", h.CreateBannerHandler)
		banner.GET("", h.GetBannersHandler)
		banner.PATCH("/:id", h.UpdateBannerHandler)
		banner.DELETE("/:id", h.DeleteBannerHandler)
	}

	router.GET("/user_banner", h.userIdentity, h.GetUserBannerHandler)
	return router
}
