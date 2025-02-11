package handlers

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandlers struct {
	authService service.AuthService
}

func NewAuthHandlers(authService service.AuthService) *AuthHandlers {
	return &AuthHandlers{authService: authService}
}

// RegisterHandler Регистрация пользователя
func (h *AuthHandlers) RegisterHandler(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.authService.CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// LoginHandler Логин пользователя
func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, token, err := h.authService.AuthenticateUser(c, credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"user_id": user.ID,
		"token":   token,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserBalanceHandler Получение баланса пользователя
func (h *AuthHandlers) GetUserBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserBalance(c, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"user_id": user.ID,
		"balance": user.Coins,
	}

	c.JSON(http.StatusOK, response)
}
