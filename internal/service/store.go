package service

import (
	"avito_shop/internal/repository"
	"context"
	"errors"
	"fmt"
)

type StoreService interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}

type storeService struct {
	userRepo repository.UserRepository
	itemRepo repository.ItemRepository
}

func NewStoreService(userRepo repository.UserRepository, itemRepo repository.ItemRepository) StoreService {
	return &storeService{userRepo: userRepo, itemRepo: itemRepo}
}

// BuyItem осуществляет покупку предмета
func (s *storeService) BuyItem(ctx context.Context, userID int, itemName string) error {
	// Поиск предмета по имени
	item, err := s.itemRepo.GetItemByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("товар не найден")
	}

	// Поиск пользователя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	// Проверяем баланс
	if user.Balance < item.Price {
		return errors.New("недостаточно средств")
	}

	// Списываем средства
	newBalance := user.Balance - item.Price
	err = s.userRepo.UpdateUserBalance(ctx, userID, newBalance)
	if err != nil {
		return fmt.Errorf("не удалось обновить баланс: %w", err)
	}

	// Добавляем предмет в инвентарь пользователя
	if err := s.itemRepo.AddToInventory(ctx, userID, item.ID); err != nil {
		return fmt.Errorf("не удалось добавить предмет в инвентарь: %w", err)
	}

	return nil
}
