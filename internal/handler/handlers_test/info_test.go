package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito_shop/internal/handler"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInfoService мок для InfoService.
type MockInfoService struct {
	mock.Mock
}

func (m *MockInfoService) GetUserInfo(ctx context.Context, userID int) (*service.UserInfoResponse, error) {
	args := m.Called(ctx, userID)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	userInfo, ok := args.Get(0).(*service.UserInfoResponse)
	if !ok {
		return nil, errors.New("failed to cast to *service.UserInfoResponse")
	}

	return userInfo, nil
}
func TestGetUserInfoHandler(t *testing.T) {
	mockInfoService := new(MockInfoService)
	infoHandler := handler.NewInfoHandler(mockInfoService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/info", infoHandler.GetUserInfoHandler)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "пользователь не аутентифицирован",
			setupRequest: func() *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/info", nil)
				if err != nil {
					t.Fatalf("Ошибка при создании запроса: %v", err)
				}
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"пользователь не аутентифицирован"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
