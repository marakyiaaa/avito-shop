package handlers_test

//
//import (
//	"avito_shop/internal/handlers"
//	"avito_shop/internal/models/entities"
//	"context"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//type MockAuthService struct {
//	mock.Mock
//}
//
//func (m *MockAuthService) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error) {
//	args := m.Called(ctx, username, password)
//	return args.Get(0).(*entities.User), args.String(1), args.Error(2)
//}
//func TestAuthHandler(t *testing.T) {
//	mockAuthService := new(MockAuthService)
//	authHandler := handlers.NewAuthHandlers(mockAuthService)
//
//	gin.SetMode(gin.TestMode)
//	router := gin.Default()
//	router.POST("/auth", authHandler.AuthHandler)
//
//	t.Run("успешная аутентификация", func(t *testing.T) {
//		mockAuthService.On("AuthenticateUser", mock.Anything, "testuser", "password123").
//			Return(&entities.User{ID: 1}, "token123", nil)
//
//		reqBody := `{"username": "testuser", "password": "password123"}`
//		req, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
//		req.Header.Set("Content-Type", "application/json")
//
//		w := httptest.NewRecorder()
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//
//		expectedBody := `{"user_id":1,"token":"token123"}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//
//		mockAuthService.AssertExpectations(t)
//	})
//
//	t.Run("ошибка валидации запроса", func(t *testing.T) {
//		reqBody := `{"username": "", "password": ""}`
//		req, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
//		req.Header.Set("Content-Type", "application/json")
//
//		w := httptest.NewRecorder()
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//
//		expectedBody := `{"errors":"Key: 'AuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag; Key: 'AuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//	})
//
//	t.Run("ошибка аутентификации", func(t *testing.T) {
//		mockAuthService.On("AuthenticateUser", mock.Anything, "wronguser", "wrongpass").
//			Return((*entities.User)(nil), "", errors.New("invalid credentials"))
//
//		reqBody := `{"username": "wronguser", "password": "wrongpass"}`
//		req, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
//		req.Header.Set("Content-Type", "application/json")
//
//		w := httptest.NewRecorder()
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusUnauthorized, w.Code)
//
//		expectedBody := `{"errors":"invalid credentials"}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//
//		mockAuthService.AssertExpectations(t)
//	})
//}
