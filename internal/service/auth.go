package service

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService interface {
	AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error)
	GenerateJWT(user *entities.User) (string, error)
	//GetUserBalance(ctx context.Context, userID int) (*entities.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	secretKey string
}

func NewAuthService(userRepo repository.UserRepository, secretKey string) AuthService {
	return &authService{userRepo: userRepo, secretKey: secretKey}
}

// AuthenticateUser Проверка пользователя, регистрация при необходимости и генерация токена
func (s *authService) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error) {
	// Получаем пользователя из БД
	user, err := s.userRepo.GetUserByUsername(ctx, username)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ошибка при получении пользователя: %v", err) // Логируем ошибку
		return nil, "", errors.New("ошибка аутентификации: invalid username or password")
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

		// После создания пользователя, снова получаем его из базы данных
		user, err = s.userRepo.GetUserByUsername(ctx, username)
		if err != nil {
			log.Printf("Ошибка при получении пользователя после регистрации: %v", err)
			return nil, "", errors.New("ошибка аутентификации: failed to retrieve user after registration")
		}
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Генерируем JWT
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// GenerateJWT Генерация JWT токена
func (s *authService) GenerateJWT(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		log.Printf("Ошибка генерации JWT: %v", err)
		return "", errors.New("failed to generate token")
	}
	return tokenString, nil
}

// CreateUser Регистрация пользователя
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
		return errors.New("failed to create user")
	}
	return nil
}

func (s *authService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return "", errors.New("failed to hash password")
	}
	return string(hash), nil
}
