package handlers

//import (
//	"avito_shop/internal/service"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strconv"
//)
//
//type TransactionHandler struct {
//	transactionService *service.TransactionService
//}
//
//func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
//	return &TransactionHandler{transactionService: service}
//}
//
//// Отправка монет
//func (h *TransactionHandler) SendCoinHandler(c *gin.Context) {
//	// Получаем userID из контекста
//	userID, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
//		return
//	}
//
//	// Получаем recipientID и amount
//	recipientIDStr := c.Param("recipient_id")
//	recipientID, err := strconv.Atoi(recipientIDStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
//		return
//	}
//
//	amountStr := c.DefaultQuery("amount", "0")
//	amount, err := strconv.ParseFloat(amountStr, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
//		return
//	}
//
//	// Отправляем монеты
//	err = h.transactionService.SendCoins(userID.(int), recipientID, amount)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Coins sent successfully"})
//}

//// Обработчик для отправки монет
//func (h *StoreHandler) SendCoinHandler(c *gin.Context) {
//	// Получаем userID из контекста
//	userID, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Unauthorized"})
//		return
//	}
//
//	// Получаем recipientID и сумму из параметров запроса
//	recipientID := c.Param("recipient_id")
//	amount := c.DefaultQuery("amount", "0")
//
//	// Преобразуем amount в нужный формат
//	amountFloat, err := strconv.ParseFloat(amount, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "Invalid amount"})
//		return
//	}
//
//	// Вызываем метод SendCoin из сервиса
//	err = h.service.SendCoin(c, userID.(int), recipientID, int(amountFloat))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: err.Error()})
//		return
//	}
//
//	// Ответ успешного выполнения
//	c.JSON(http.StatusOK, response.SuccessResponse{Message: "Coins sent successfully"})
//}
