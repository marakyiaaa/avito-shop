package middleware

import (
	"avito_shop/internal/models/api/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

// GenerateJWT генерирует JWT для заданного идентификатора пользователя.
func GenerateJWT(secretKey string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Ошибка генерации JWT: %v", err)
		return "", fmt.Errorf("failed to generate token")
	}
	return tokenString, nil
}

// NewCheckAuth создает новую промежуточную функцию для проверки аутентификации пользователя.
// Она извлекает JWT-токен из заголовка Authorization, проверяет его и добавляет идентификатор пользователя в контекст Gin.
func NewCheckAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid or expired token"})
			c.Abort()
			return
		}

		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid or expired token - 2"})
			c.Abort()
			return
		}
		c.Set("user_id", int(userID))
		c.Next()
	}
}
