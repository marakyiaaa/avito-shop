package service_test

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/service"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// TestAuthService_AuthenticateUser_Success Успешная аутентификация
func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Ошибка при хэшировании пароля: %v", err)
	}

	expectedUser := &entities.User{
		ID:       1,
		Username: "user1",
		Password: string(hashedPassword),
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return(expectedUser, nil)

	authService := service.NewAuthService(mockUserRepo, "secretkey")
	user, token, err := authService.AuthenticateUser(context.Background(), "user1", "password")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_AuthenticateUser_InvalidPassword Ввод неверного пароля
func TestAuthService_AuthenticateUser_InvalidPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Ошибка при хэшировании пароля: %v", err)
	}

	expectedUser := &entities.User{
		ID:       1,
		Username: "user1",
		Password: string(hashedPassword),
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return(expectedUser, nil)

	authService := service.NewAuthService(mockUserRepo, "secretkey")
	user, token, err := authService.AuthenticateUser(context.Background(), "user1", "password1")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_AuthenticateUser_UserNotFound Пользователь не найден и поседующая регистрация
func TestAuthService_AuthenticateUser_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Ошибка при хэшировании пароля: %v", err)
	}

	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return((*entities.User)(nil), sql.ErrNoRows).Once()
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)

	expectedUser := &entities.User{
		ID:       1,
		Username: "user1",
		Password: string(hashedPassword),
		Balance:  1000,
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return(expectedUser, nil).Once()

	authService := service.NewAuthService(mockUserRepo, "secretkey")
	user, token, err := authService.AuthenticateUser(context.Background(), "user1", "password")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_AuthenticateUser_DatabaseErrorПри ошибке базы данных возвращается соответствующая ошибка
func TestAuthService_AuthenticateUser_DatabaseError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return((*entities.User)(nil), errors.New("database error"))

	authService := service.NewAuthService(mockUserRepo, "secretkey")
	user, token, err := authService.AuthenticateUser(context.Background(), "user1", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_AuthenticateUser_TokenGenerationError Ошибка генерации токена
func TestAuthService_AuthenticateUser_TokenGenerationError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Ошибка при хэшировании пароля: %v", err)
	}

	expectedUser := &entities.User{
		ID:       1,
		Username: "user1",
		Password: string(hashedPassword),
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, "user1").Return(expectedUser, nil)
	authService := service.NewAuthService(mockUserRepo, "") // Пустой ключ вызовет ошибку генерации токена

	user, token, err := authService.AuthenticateUser(context.Background(), "user1", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_AuthenticateUser_CreateUserError Ошибка при создании пользователя
func TestAuthService_AuthenticateUser_CreateUserError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	mockUserRepo.On("GetUserByUsername", mock.Anything, "newuser").Return((*entities.User)(nil), sql.ErrNoRows)
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).Return(errors.New("create user error"))

	authService := service.NewAuthService(mockUserRepo, "secretkey")
	user, token, err := authService.AuthenticateUser(context.Background(), "newuser", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}
