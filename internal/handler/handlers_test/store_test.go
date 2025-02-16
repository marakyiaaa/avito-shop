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

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "пользователь не аутентифицирован",
			method:         http.MethodPost,
			url:            "/buy/cup",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"неавторизованный"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
