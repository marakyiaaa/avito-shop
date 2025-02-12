package middleware

import (
	"avito_shop/internal/models/api/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

//хз новое

// NewCheckAuth — Middleware для проверки JWT
func NewCheckAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
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
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Errors: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Добавляем user_id в контекст Gin
		c.Set("user_id", int(userID))

		// Продолжаем выполнение запроса
		c.Next()
	}
}

//type Middleware func(http.HandlerFunc) http.HandlerFunc
//
//// NewCheckAuth  JWT
//func NewCheckAuth(secretKey string) Middleware {
//	return func(h http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			// Получаем токен из заголовков
//			tokenString := r.Header.Get("Authorization")
//			if tokenString == "" {
//				http.Error(w, "Token is missing", http.StatusUnauthorized)
//				return
//			}
//
//			// Парсим токен
//			parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//				// Проверяем, что метод подписи правильный
//				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//					log.Println("Unexpected signing method")
//					return nil, errors.New("invalid token")
//				}
//				return []byte(secretKey), nil
//			})
//
//			if err != nil || !parsedToken.Valid {
//				http.Error(w, "Invalid token", http.StatusUnauthorized)
//				return
//			}
//
//			// Извлекаем данные из токена
//			claims, ok := parsedToken.Claims.(jwt.MapClaims)
//			if !ok {
//				http.Error(w, "Invalid token", http.StatusUnauthorized)
//				return
//			}
//
//			// Извлекаем информацию о пользователе
//			userID := claims["user_id"].(float64)
//
//			// Передаем данные о пользователе в контекст запроса
//			ctx := context.WithValue(r.Context(), "user_id", int(userID))
//
//			h(w, r.WithContext(ctx))
//		}
//	}
//}

////ERROR GIN
//func NewCheckAuth(secretKey string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Получаем токен из заголовка Authorization
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
//			c.Abort()
//			return
//		}
//
//		// Разбиваем заголовок: "Bearer <token>"
//		tokenParts := strings.Split(authHeader, " ")
//		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
//			c.Abort()
//			return
//		}
//
//		tokenString := tokenParts[1]
//
//		// Разбираем и проверяем токен
//		claims := &jwt.MapClaims{}
//		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexpected signing method")
//			}
//			return []byte(secretKey), nil
//		})
//
//		if err != nil || !token.Valid {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
//			c.Abort()
//			return
//		}
//
//		// Получаем user_id из токена
//		userID, ok := (*claims)["user_id"].(float64)
//		if !ok {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
//			c.Abort()
//			return
//		}
//
//		// Добавляем user_id в контекст Gin
//		c.Set("user_id", int(userID))
//
//		// Продолжаем выполнение запроса
//		c.Next()
//	}
//}
