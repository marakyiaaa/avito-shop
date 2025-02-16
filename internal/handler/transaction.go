package handler

import (
	"net/http"
	"strconv"

	"avito_shop/internal/models/api/request"
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "невалидный запрос"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "сумма должна быть положительной"})
		return
	}

	fromUserID := c.GetInt("user_id")
	if fromUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "пользователь не аутентифицирован"})
		return
	}

	toUserID, err := strconv.Atoi(req.ToUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Errors: "недействительный получатель"})
		return
	}

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
