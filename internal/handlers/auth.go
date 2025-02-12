package handlers

import (
	"avito_shop/internal/models/api/request"
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthHandlers struct {
	authService service.AuthService
}

func NewAuthHandlers(authService service.AuthService) *AuthHandlers {
	return &AuthHandlers{authService: authService}
}

func (h *AuthHandlers) AuthHandler(c *gin.Context) {
	var authReq request.AuthRequest

	if err := c.ShouldBindJSON(&authReq); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(ve)})
			return
		}
	}

	user, token, err := h.authService.AuthenticateUser(c, authReq.Username, authReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp := response.AuthResponse{Token: token}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"token":   resp.Token,
	})
}

func formatValidationErrors(ve validator.ValidationErrors) []string {
	var errors []string
	for _, fe := range ve {
		errors = append(errors, fe.Error()) // Использует встроенное описание ошибки
	}
	return errors
}
