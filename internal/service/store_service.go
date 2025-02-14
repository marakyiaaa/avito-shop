package service

import (
	"avito_shop/internal/repository"
	"context"
	"fmt"
	"log"
)

// StoreService определяет методы для работы с магазином.
type StoreService interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}

// storeService реализует интерфейс StoreService.
// Использует репозитории пользователей и предметов и инвентаря для выполнения операций.
type storeService struct {
	userRepo      repository.UserRepository
	itemRepo      repository.ItemRepository
	inventoryRepo repository.InventoryRepository
}

// NewStoreService создает новый экземпляр storeService.
// Принимает репозитории пользователей и предметов и инвентаря.
func NewStoreService(userRepo repository.UserRepository, itemRepo repository.ItemRepository, inventoryRepo repository.InventoryRepository) StoreService {
	return &storeService{userRepo: userRepo, itemRepo: itemRepo, inventoryRepo: inventoryRepo}
}

// BuyItem осуществляет покупку предмета пользователем.
// Проверяет баланс пользователя и добавляет предмет в его инвентарь.
func (s *storeService) BuyItem(ctx context.Context, userID int, itemName string) error {
	item, err := s.itemRepo.GetItemByName(ctx, itemName)
	if err != nil || item == nil {
		log.Println("Ошибка при получении товара")
		return fmt.Errorf("товар не найден")
	}
	log.Println("Цена товара:", item.Price)

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		log.Println("Ошибка при получении пользователя:", err)
		return fmt.Errorf("пользователь не найден")
	}

	if user.Balance < item.Price {
		log.Println("Недостаточно средств")
		return fmt.Errorf("недостаточно средств")
	}

	newBalance := user.Balance - item.Price
	err = s.userRepo.UpdateUserBalance(ctx, userID, newBalance)
	if err != nil {
		log.Println("Ошибка при обновлении баланса:", err)
		return fmt.Errorf("не удалось обновить баланс: %w", err)
	}

	if err := s.inventoryRepo.AddToInventory(ctx, userID, itemName); err != nil {
		return fmt.Errorf("не удалось добавить предмет в инвентарь: %w", err)
	}

	log.Println("Покупка завершена успешно")
	return nil
}
