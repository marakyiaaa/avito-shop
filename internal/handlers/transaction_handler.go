package handlers

import (
	"avito_shop/internal/models/api/request"
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TransactionHandler - транзакции между пользователями
type TransactionHandler struct {
	transactionService service.TransactionService
}

// NewTransactionHandler конструктор (создает новый экземпляр) TransactionHandler.
func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// SendCoinHandler - обработчик для отправки монет
func (h *TransactionHandler) SendCoinHandler(c *gin.Context) {
	var req request.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "invalid request"})
		return
	}

	// Добавляем проверку на отрицательное значение
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "amount must be a positive value"})
		return
	}

	// Получаем ID пользователя из токена
	fromUserID := c.GetInt("user_id")
	if fromUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "user not authenticated"})
		return
	}

	// Преобразуем имя получателя в ID (если нужно, запросить у репозитория)
	toUserID, err := strconv.Atoi(req.ToUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "invalid recipient"})
		return
	}

	// Отправляем монеты
	err = h.transactionService.SendCoins(c, fromUserID, toUserID, req.Amount)
	if err != nil {
		if err.Error() == "нельзя отправлять монеты самому себе" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Errors: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
