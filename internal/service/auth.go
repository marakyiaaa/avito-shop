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
	CreateUser(ctx context.Context, user *entities.User) error
	AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error)
	GenerateJWT(user *entities.User) (string, error)

	GetUserBalance(ctx context.Context, userID int) (*entities.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	secretKey string
}

func NewAuthService(userRepo repository.UserRepository, secretKey string) AuthService {
	return &authService{userRepo: userRepo, secretKey: secretKey}
}

// CreateUser Регистрация пользователя
func (s *authService) CreateUser(ctx context.Context, user *entities.User) error {
	//Существует ли username
	existUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existUser != nil {
		return errors.New("username already taken")
	}

	//хэшируем пароль
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return errors.New("failed to create user")
	}
	//Записываем хэшированный пароль
	user.Password = string(hashPassword)

	/// Создаем пользователя в БД
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		return errors.New("failed to create user")
	}
	return nil
}

// AuthenticateUser Проверка пользователя и создание токена
func (s *authService) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, string, error) {
	// Получаем пользователя из БД
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)   // Логируем ошибку
		return nil, "", errors.New("invalid username or password") // но не раскрываем детали пользователю
	}

	//// Если пользователь не найден
	//if user == nil {
	//	return nil, errors.New("invalid username or password") // преоверка в GetUserByUsername есть
	//}

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
		return "", err
	}
	return tokenString, nil
}

func (s *authService) GetUserBalance(ctx context.Context, userID int) (*entities.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
