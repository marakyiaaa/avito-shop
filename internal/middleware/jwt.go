package middleware

import (
	"avito_shop/internal/models/api/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

func GenerateJWT(secretKey string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Ошибка генерации JWT: %v", err)
		return "", errors.New("failed to generate token")
	}
	return tokenString, nil
}

func NewCheckAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Missing Authorization header"})
			c.Abort()
			return
		}

		// Разбиваем заголовок: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Разбираем и проверяем токен
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

		// Получаем user_id из токена
		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid or expired token - 2"})
			c.Abort()
			return
		}

		// Добавляем user_id в контекст Gin
		c.Set("user_id", int(userID))

		// Продолжаем выполнение запроса
		c.Next()
	}
}
