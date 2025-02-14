package service

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"fmt"
)

type InfoService interface {
	GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error)
}

type UserInfoResponse struct {
	User         *entities.User
	Inventory    []*entities.Inventory
	Transactions []*entities.Transaction
}

type infoService struct {
	userRepo        repository.UserRepository
	itemRepo        repository.ItemRepository
	transactionRepo repository.TransactionRepository
}

func NewInfoService(userRepo repository.UserRepository, itemRepo repository.ItemRepository, transactionRepo repository.TransactionRepository) InfoService {
	return &infoService{userRepo: userRepo, itemRepo: itemRepo, transactionRepo: transactionRepo}
}

func (s *infoService) GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	inventory, err := s.itemRepo.GetInventoryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении инвентаря: %w", err)
	}

	transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении транзакций: %w", err)
	}

	return &UserInfoResponse{
		User:         user,
		Inventory:    inventory,
		Transactions: transactions,
	}, nil
}
