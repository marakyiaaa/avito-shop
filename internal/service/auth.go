package service

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

//начало

type AuthService interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func newAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// CreateUser Регистрация пользователя
func (s *authService) CreateUser(ctx context.Context, user *entities.User) error {
	//существует ли username
	existUser, _ := s.userRepo.GetUserByUsername(ctx, user.Username)
	if existUser != nil {
		return errors.New("username already taken")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Username), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashPassword)

	return s.userRepo.CreateUser(ctx, user)
}

//
//func (s *authService) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
//	return s.userRepo.GetUserByID(ctx, id)
//}
