package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito_shop/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStoreService мок для StoreService.
type MockStoreService struct {
	mock.Mock
}

func (m *MockStoreService) BuyItem(ctx context.Context, userID int, itemName string) error {
	args := m.Called(ctx, userID, itemName)
	return args.Error(0)
}

func TestBuyItemHandler(t *testing.T) {
	mockStoreService := new(MockStoreService)

	storeHandler := handler.NewStoreHandler(mockStoreService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/buy/:item", storeHandler.BuyItemHandler)

	t.Run("пользователь не аутентифицирован", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/buy/cup", nil)
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		expectedBody := `{"errors":"неавторизованный"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})
}
