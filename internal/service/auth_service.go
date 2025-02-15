package service

import (
	"avito_shop/internal/middleware"
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	// Получаем пользователя из БД
	user, err := s.userRepo.GetUserByUsername(ctx, username)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ошибка при получении пользователя: %v", err) // Логируем ошибку
		return nil, "", fmt.Errorf("ошибка аутентификации: invalid username or password")
	}

	// Если пользователь не найден — регистрируем его
	if user == nil {
		newUser := &entities.User{
			Username: username,
			Password: password,
			Balance:  1000,
		}

		if err := s.createUser(ctx, newUser); err != nil {
			return nil, "", err
		}

		// После создания пользователя, снова получаем его из БД
		user, err = s.userRepo.GetUserByUsername(ctx, username)
		if err != nil {
			log.Printf("Ошибка при получении пользователя после регистрации: %v", err)
			return nil, "", fmt.Errorf("ошибка аутентификации: failed to retrieve user after registration")
		}
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid username or password")
	}

	// Генерируем JWT
	token, err := middleware.GenerateJWT(s.secretKey, user.ID)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return nil, "", fmt.Errorf("ошибка аутентификации: failed to generate token")
	}
	return user, token, nil
}

// createUser регистрирует нового пользователя.
func (s *authService) createUser(ctx context.Context, user *entities.User) error {
	// Хэшируем пароль
	hashPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}
	//Записываем хэшированный пароль
	user.Password = hashPassword

	// Создаем пользователя в БД
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		return fmt.Errorf("failed to create user")
	}
	return nil
}

// hashPassword хэширует пароль пользователя.
func (s *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return "", fmt.Errorf("failed to hash password")
	}
	return string(hash), nil
}
