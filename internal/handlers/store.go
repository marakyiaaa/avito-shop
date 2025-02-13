package handlers

import (
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(service service.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

// BuyItemHandler обрабатывает покупку предмета
func (h *StoreHandler) BuyItemHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Unauthorized"})
		return
	}

	// Получаем название товара
	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "Не указано название товара"})
		return
	}

	// Совершаем покупку через сервис
	err := h.service.BuyItem(c, userID.(int), itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Transaction failed"})
		return
	}
}
