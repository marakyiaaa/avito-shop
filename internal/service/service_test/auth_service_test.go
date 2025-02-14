package service_test

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/service"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// TestAuthService_AuthenticateUser_Success — тест для успешной аутентификации.
func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	// Создаем мок для UserRepository
	mockUserRepo := new(MockUserRepository)

	// Настраиваем мок
	expectedUser := &entities.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$examplehash", // Хэшированный пароль
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(expectedUser, nil)

	// Создаем сервис с моком
	authService := service.NewAuthService(mockUserRepo, "secretkey")

	// Вызываем метод
	user, token, err := authService.AuthenticateUser(context.Background(), "testuser", "password")

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_InvalidPassword(t *testing.T) {
	// Создаем мок для UserRepository
	mockUserRepo := new(MockUserRepository)

	// Настраиваем мок
	expectedUser := &entities.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$examplehash", // Хэшированный пароль
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(expectedUser, nil)

	// Создаем сервис с моком
	authService := service.NewAuthService(mockUserRepo, "secretkey")

	// Вызываем метод с неверным паролем
	user, token, err := authService.AuthenticateUser(context.Background(), "testuser", "wrongpassword")

	// Проверяем результаты
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}
