package handlers_test

//
//import (
//	"avito_shop/internal/handlers"
//	"avito_shop/internal/models/entities"
//	"avito_shop/internal/service"
//	"context"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//// MockInfoService мок для InfoService.
//type MockInfoService struct {
//	mock.Mock
//}
//
//func (m *MockInfoService) GetUserInfo(ctx context.Context, userID int) (*service.UserInfoResponse, error) {
//	args := m.Called(ctx, userID)
//	return args.Get(0).(*service.UserInfoResponse), args.Error(1)
//}
//
//func TestGetUserInfoHandler(t *testing.T) {
//	// Создаем мок сервиса
//	mockInfoService := new(MockInfoService)
//
//	// Создаем обработчик с моком сервиса
//	infoHandler := handlers.NewInfoHandler(mockInfoService)
//
//	// Настраиваем Gin для тестирования
//	gin.SetMode(gin.TestMode)
//	router := gin.Default()
//	router.GET("/info", infoHandler.GetUserInfoHandler)
//
//	t.Run("успешное получение информации о пользователе", func(t *testing.T) {
//		// Ожидаем вызов метода GetUserInfo с определенными аргументами
//		mockInfoService.On("GetUserInfo", mock.Anything, 1).
//			Return(&service.UserInfoResponse{
//				User: &entities.User{
//					ID:       1,
//					Username: "testuser",
//					Balance:  1000,
//				},
//				Inventory: []*entities.Inventory{
//					{ItemType: "book", Quantity: 1},
//				},
//				Transactions: []*entities.Transaction{
//					{FromUserID: 2, ToUserID: 1, Amount: 100},
//				},
//			}, nil)
//
//		// Создаем тестовый запрос
//		req, _ := http.NewRequest(http.MethodGet, "/info", nil)
//
//		// Устанавливаем user_id в контекст Gin
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = req
//		c.Set("user_id", 1) // Устанавливаем user_id в контекст
//
//		// Вызываем обработчик
//		router.ServeHTTP(w, req)
//
//		// Проверяем статус код
//		assert.Equal(t, http.StatusOK, w.Code)
//
//		// Проверяем тело ответа
//		expectedBody := `{"user":{"id":1,"username":"testuser","balance":1000},"inventory":[{"item_type":"book","quantity":1}],"transactions":[{"from_user_id":2,"to_user_id":1,"amount":100}]}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//
//		// Проверяем, что все ожидания выполнены
//		mockInfoService.AssertExpectations(t)
//	})
//
//	t.Run("пользователь не аутентифицирован", func(t *testing.T) {
//		// Создаем тестовый запрос
//		req, _ := http.NewRequest(http.MethodGet, "/info", nil)
//
//		// Устанавливаем пустой контекст (без user_id)
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = req
//
//		// Вызываем обработчик
//		router.ServeHTTP(w, req)
//
//		// Проверяем статус код
//		assert.Equal(t, http.StatusUnauthorized, w.Code)
//
//		// Проверяем тело ответа
//		expectedBody := `{"errors":"user not authenticated"}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//	})
//
//	t.Run("ошибка при получении информации о пользователе", func(t *testing.T) {
//		// Ожидаем вызов метода GetUserInfo с определенными аргументами
//		mockInfoService.On("GetUserInfo", mock.Anything, 1).
//			Return((*service.UserInfoResponse)(nil), errors.New("database error"))
//
//		// Создаем тестовый запрос
//		req, _ := http.NewRequest(http.MethodGet, "/info", nil)
//
//		// Устанавливаем user_id в контекст Gin
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = req
//		c.Set("user_id", 1) // Устанавливаем user_id в контекст
//
//		// Вызываем обработчик
//		router.ServeHTTP(w, req)
//
//		// Проверяем статус код
//		assert.Equal(t, http.StatusInternalServerError, w.Code)
//
//		// Проверяем тело ответа
//		expectedBody := `{"errors":"database error"}`
//		assert.JSONEq(t, expectedBody, w.Body.String())
//
//		// Проверяем, что все ожидания выполнены
//		mockInfoService.AssertExpectations(t)
//	})
//}
