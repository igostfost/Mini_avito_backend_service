package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}
	// parse token

	userId, isAdmin, err := h.utils.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("userId", userId)
	c.Set("isAdmin", isAdmin)
}
