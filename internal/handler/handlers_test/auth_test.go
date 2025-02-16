package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"avito_shop/internal/handler"
	"avito_shop/internal/models/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error) {
	args := m.Called(ctx, username, password)

	user, ok := args.Get(0).(*entities.User)
	if !ok && args.Get(0) != nil {
		panic("Ошибка приведения типа *entities.User")
	}

	return user, args.String(1), args.Error(2)
}
func TestAuthHandler(t *testing.T) {
	mockAuthService := new(MockAuthService)
	authHandler := handler.NewAuthHandlers(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth", authHandler.AuthHandler)

	t.Run("успешная аутентификация", func(t *testing.T) {
		mockAuthService.On("AuthenticateUser", mock.Anything, "testuser", "password123").
			Return(&entities.User{ID: 1}, "token123", nil)

		reqBody := `{"username": "testuser", "password": "password123"}`
		req, err := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedBody := `{"user_id":1,"token":"token123"}`
		assert.JSONEq(t, expectedBody, w.Body.String())

		mockAuthService.AssertExpectations(t)
	})
	t.Run("ошибка валидации запроса", func(t *testing.T) {
		reqBody := `{"username": "", "password": ""}`
		req, err := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := `{"errors":"Key: 'AuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag; Key: 'AuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

	t.Run("ошибка аутентификации", func(t *testing.T) {
		mockAuthService.On("AuthenticateUser", mock.Anything, "wronguser", "wrongpass").
			Return((*entities.User)(nil), "", errors.New("invalid credentials"))

		reqBody := `{"username": "", "password": ""}`
		req, err := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedBody := `{"errors":"Key: 'AuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag; Key: 'AuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})
}
