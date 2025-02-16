package service

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"fmt"
)

// InfoService определяет методы для получения информации о пользователе.
type InfoService interface {
	GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error)
}

// UserInfoResponse представляет ответ с информацией о пользователе.
// Включает данные пользователя, его инвентарь и список транзакций.
type UserInfoResponse struct {
	User         *entities.User
	Inventory    []*entities.Inventory
	Transactions []*entities.Transaction
}

// infoService реализует интерфейс InfoService.
// Использует репозитории пользователей, предметов, инвентаря и транзакций для получения данных.
type infoService struct {
	userRepo        repository.UserRepository
	itemRepo        repository.ItemRepository
	inventoryRepo   repository.InventoryRepository
	transactionRepo repository.TransactionRepository
}

// NewInfoService создает новый экземпляр infoService.
// Принимает репозитории пользователей, предметов, инвентаря и транзакций.
func NewInfoService(userRepo repository.UserRepository, itemRepo repository.ItemRepository, inventoryRepo repository.InventoryRepository, transactionRepo repository.TransactionRepository) InfoService {
	return &infoService{userRepo: userRepo, itemRepo: itemRepo, inventoryRepo: inventoryRepo, transactionRepo: transactionRepo}
}

// GetUserInfo возвращает информацию о пользователе, включая его инвентарь и транзакции.
func (s *infoService) GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	inventory, err := s.inventoryRepo.GetInventoryByUserID(ctx, userID)
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
