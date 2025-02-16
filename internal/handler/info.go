package handler

import (
	"net/http"

	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
)

// InfoHandler - получение информации пользователя:
// баланс, список купленных им товаров, транзакции между пользователями.
type InfoHandler struct {
	service service.InfoService
}

// NewInfoHandler конструктор (создает новый экземпляр) InfoHandler.
func NewInfoHandler(service service.InfoService) *InfoHandler {
	return &InfoHandler{service: service}
}

// GetUserInfoHandler обрабатывает запрос на получение информации о пользователе,
// отправляет ответ обратно клиенту.
func (h *InfoHandler) GetUserInfoHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "пользователь не аутентифицирован"})
		return
	}

	info, err := h.service.GetUserInfo(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
