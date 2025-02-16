package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"avito_shop/internal/middleware"
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет методы для аутентификации пользователей.
type AuthService interface {
	AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error)
}

// authService реализует интерфейс AuthService.
// Использует репозиторий пользователей для работы с данными.
type authService struct {
	userRepo  repository.UserRepository
	secretKey string
}

// NewAuthService создает новый экземпляр authService.
// Принимает репозиторий пользователей и секретный ключ для генерации JWT.
func NewAuthService(userRepo repository.UserRepository, secretKey string) AuthService {
	return &authService{userRepo: userRepo, secretKey: secretKey}
}

// AuthenticateUser проверяет учетные данные пользователя и генерирует JWT-токен.
func (s *authService) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return nil, "", fmt.Errorf("ошибка аутентификации")
	}

	if user == nil {
		newUser := &entities.User{
			Username: username,
			Password: password,
			Balance:  1000,
		}

		if err := s.createUser(ctx, newUser); err != nil {
			return nil, "", err
		}

		user, err = s.userRepo.GetUserByUsername(ctx, username)
		if err != nil {
			log.Printf("Ошибка при получении пользователя после регистрации: %v", err)
			return nil, "", fmt.Errorf("ошибка аутентификации: при получении пользователя после регистрации")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("ошибка аутентификации: неверный пароль")
	}

	token, err := middleware.GenerateJWT(s.secretKey, user.ID)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return nil, "", fmt.Errorf("ошибка аутентификации: токен не сгенерировался")
	}
	return user, token, nil
}

// createUser регистрирует нового пользователя.
func (s *authService) createUser(ctx context.Context, user *entities.User) error {
	hashPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashPassword

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Printf("ошибка создания пользователя: %v", err)
		return fmt.Errorf("ошибка создания пользователя")
	}
	return nil
}

// hashPassword хэширует пароль пользователя.
func (s *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ошибка хеширования пароля: %v", err)
		return "", fmt.Errorf("ошибка хеширования пароля")
	}
	return string(hash), nil
}
