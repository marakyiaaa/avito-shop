package handler

import (
	"net/http"

	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"

	"github.com/gin-gonic/gin"
)

// StoreHandler - покупка товаров в магазине.
type StoreHandler struct {
	service service.StoreService
}

// NewStoreHandler конструктор (создает новый экземпляр) StoreHandler.
func NewStoreHandler(service service.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

// BuyItemHandler обрабатывает покупку предмета
func (h *StoreHandler) BuyItemHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "неавторизованный"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Ошибка сервера: неверный формат user_id"})
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "Не указано название товара"})
		return
	}

	err := h.service.BuyItem(c, id, itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Сбой транзакции"})
		return
	}
}
