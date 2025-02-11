package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

//хз новое

type Middleware func(http.HandlerFunc) http.HandlerFunc

// NewCheckAuth Проверка JWT
func NewCheckAuth(secretKey string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовков
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Token is missing", http.StatusUnauthorized)
				return
			}

			// Парсим токен
			parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что метод подписи правильный
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Println("Unexpected signing method")
					return nil, errors.New("invalid token")
				}
				return []byte(secretKey), nil
			})

			if err != nil || !parsedToken.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Извлекаем данные из токена
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Извлекаем информацию о пользователе
			userID := claims["user_id"].(float64)

			// Передаем данные о пользователе в контекст запроса
			ctx := context.WithValue(r.Context(), "user_id", int(userID))

			h(w, r.WithContext(ctx))
		}
	}
}
