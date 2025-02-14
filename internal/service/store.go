package service

import (
	"avito_shop/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
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
		log.Println("Ошибка при получении товара:", err)
		return fmt.Errorf("товар не найден")
	}
	log.Println("Цена товара:", item.Price)

	// Поиск пользователя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Ошибка при получении пользователя:", err)
		return fmt.Errorf("пользователь не найден: %w", err)
	}
	if user == nil {
		log.Println("Пользователь не найден")
		return errors.New("пользователь не найден")
	}

	// Проверяем баланс
	if user.Balance < item.Price {
		log.Println("Недостаточно средств")
		return errors.New("недостаточно средств")
	}

	// Списываем средства
	newBalance := user.Balance - item.Price
	err = s.userRepo.UpdateUserBalance(ctx, userID, newBalance)
	if err != nil {
		log.Println("Ошибка при обновлении баланса:", err)
		return fmt.Errorf("не удалось обновить баланс: %w", err)
	}

	// Добавляем предмет в инвентарь пользователя
	if err := s.itemRepo.AddToInventory(ctx, userID, itemName); err != nil {
		return fmt.Errorf("не удалось добавить предмет в инвентарь: %w", err)
	}

	log.Println("Покупка завершена успешно")
	return nil
}
