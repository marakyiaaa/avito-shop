package handlers

import (
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoHandler struct {
	service service.InfoService
}

func NewInfoHandler(service service.InfoService) *InfoHandler {
	return &InfoHandler{service: service}
}

func (h *InfoHandler) GetUserInfoHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	info, err := h.service.GetUserInfo(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
