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

	t.Run("пользователь не аутентифицирован", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/info", nil)
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		expectedBody := `{"errors":"пользователь не аутентифицирован"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

}
