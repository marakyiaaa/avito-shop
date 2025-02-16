package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"avito_shop/internal/models/api/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT генерирует JWT для заданного идентификатора пользователя.
func GenerateJWT(secretKey string, userID int) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("пустой ключ")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Ошибка генерации JWT: %v", err)
		return "", fmt.Errorf("ошибка генерации JWT")
	}
	return tokenString, nil
}

// NewCheckAuth создает новую промежуточную функцию для проверки аутентификации пользователя.
// Она извлекает JWT-токен из заголовка Authorization, проверяет его и добавляет идентификатор пользователя в контекст Gin.
func NewCheckAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "отсутствует заголовок "})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "неверный формат заголовка "})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("недействительный токен")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "недействительный токен"})
			c.Abort()
			return
		}

		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "недействительный токен"})
			c.Abort()
			return
		}
		c.Set("user_id", int(userID))
		c.Next()
	}
}
