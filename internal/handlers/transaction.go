package handlers

import (
	"avito_shop/internal/models/api/request"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// SendCoinHandler - обработчик для отправки монет
func (h *TransactionHandler) SendCoinHandler(c *gin.Context) {
	var req request.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Добавляем проверку на отрицательное значение
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be a positive value"})
		return
	}

	// Получаем ID пользователя из токена
	fromUserID := c.GetInt("user_id")
	if fromUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Преобразуем имя получателя в ID (если нужно, запросить у репозитория)
	toUserID, err := strconv.Atoi(req.ToUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipient"})
		return
	}

	// Отправляем монеты
	err = h.transactionService.SendCoins(c, fromUserID, toUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
