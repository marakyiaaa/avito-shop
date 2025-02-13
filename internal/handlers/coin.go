package handlers

import (
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StoreHandler struct {
	service service.StoreService
}

func NewStoreHandler(service service.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

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

	// Получаем данные о пользователе
	user, err := h.service.GetUserByID(c, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Errors: "User not found"})
		return
	}

	// Получаем информацию о предмете
	item, err := h.service.GetItemByName(c, itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Errors: "Item not found"})
		return
	}

	// Проверяем баланс пользователя
	if user.Balance < item.Price {
		c.JSON(http.StatusForbidden, response.ErrorResponse{Errors: "Not enough coins"})
		return
	}

	// Совершаем покупку
	_, err = h.service.BuyItem(c, user.ID, item.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Transaction failed"})
		return
	}

	// Обновляем данные пользователя после покупки
	updatedUser, err := h.service.GetUserByID(c, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Failed to retrieve updated user data"})
		return
	}

	// Получаем историю транзакций
	history, err := h.service.GetTransactionHistory(c, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: "Failed to retrieve transaction history"})
		return
	}

	// Формируем ответ с балансом, инвентарем и историей транзакций
	c.JSON(http.StatusOK, response.InfoResponse{
		Coins:       updatedUser.Balance,
		Inventory:   []response.InventoryItem{{Type: item.Name, Quantity: 1}},
		CoinHistory: *history,
	})

	//// Формируем ответ с актуальным балансом и инвентарем
	//c.JSON(http.StatusOK, response.InfoResponse{
	//	Coins: updatedUser.Balance,
	//	Inventory: []response.InventoryItem{
	//		{Type: item.Name, Quantity: 1}, // Можно сделать динамическое обновление инвентаря
	//	},
	//	CoinHistory: response.CoinHistory{}, // Можно добавить логику записи транзакции
	//})
}

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

// Обработчик для отправки монет
func (h *StoreHandler) SendCoinHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Unauthorized"})
		return
	}

	// Получаем recipientID и сумму из параметров запроса
	recipientID := c.Param("recipient_id")
	amount := c.DefaultQuery("amount", "0")

	// Преобразуем amount в нужный формат
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "Invalid amount"})
		return
	}

	// Вызываем метод SendCoin из сервиса
	err = h.service.SendCoin(c, userID.(int), recipientID, int(amountFloat))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: err.Error()})
		return
	}

	// Ответ успешного выполнения
	c.JSON(http.StatusOK, response.SuccessResponse{Message: "Coins sent successfully"})
}
