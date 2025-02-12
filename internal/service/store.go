package service

import (
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"errors"
	"strconv"
	"time"
)

type StoreService interface {
	GetItemByName(ctx context.Context, name string) (*entities.Item, error)
	GetUserByID(ctx context.Context, userID int) (*entities.User, error)
	GetTransactionHistory(ctx context.Context, userID int) (*response.CoinHistory, error)
	BuyItem(ctx context.Context, userID int, itemName string) (*entities.Transaction, error)
}

type storeService struct {
	userRepo  repository.UserRepository
	itemRepo  repository.ItemRepository
	transRepo repository.TransactionRepository
}

func NewStoreService(userRepo repository.UserRepository, itemRepo repository.ItemRepository, transRepo repository.TransactionRepository) StoreService {
	return &storeService{userRepo: userRepo, itemRepo: itemRepo, transRepo: transRepo}
}

// GetItemByName находит предмет по имени
func (s *storeService) GetItemByName(ctx context.Context, name string) (*entities.Item, error) {
	item, err := s.itemRepo.GetItemByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetUserByID находит пользователя по ID
func (s *storeService) GetUserByID(ctx context.Context, userID int) (*entities.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// BuyItem осуществляет покупку предмета
func (s *storeService) BuyItem(ctx context.Context, userID int, itemName string) (*entities.Transaction, error) {
	// 1. Получаем товар
	item, err := s.itemRepo.GetItemByName(ctx, itemName)
	if err != nil {
		return nil, errors.New("товар не найден")
	}

	// 2. Получаем пользователя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}
	if user == nil {
		return nil, errors.New("пользователь не существует")
	}

	// 3. Проверяем баланс
	if user.Balance < item.Price {
		return nil, errors.New("недостаточно средств")
	}

	// 4. Списываем средства
	newBalance := user.Balance - item.Price
	err = s.userRepo.UpdateUserBalance(ctx, userID, newBalance)
	if err != nil {
		return nil, err
	}

	//5. Создаём запись о покупке
	// Создаем транзакцию
	transaction := &entities.Transaction{
		FromUserID: user.ID,
		ToUserID:   0, // Система
		Amount:     item.Price,
		Timestamp:  time.Now(),
	}

	transaction, err = s.transRepo.CreateTransaction(ctx, user.ID, 0, item.Price)
	if err != nil {
		return nil, errors.New("failed to create transaction")
	}

	return transaction, nil
}

func (s *storeService) GetTransactionHistory(ctx context.Context, userID int) (*response.CoinHistory, error) {
	// Получаем отправленные транзакции (userID -> другие пользователи)
	sent, err := s.transRepo.GetSentTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем полученные транзакции (другие пользователи -> userID)
	received, err := s.transRepo.GetReceivedTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Формируем структуру ответа
	history := &response.CoinHistory{
		Sent:     make([]response.CoinTransaction, len(sent)),
		Received: make([]response.CoinTransaction, len(received)),
	}

	// Заполняем отправленные транзакции
	for i, t := range sent {
		history.Sent[i] = response.CoinTransaction{
			FromUser: strconv.Itoa(t.ToUserID), // Можно заменить на имя пользователя
			Amount:   t.Amount,
		}
	}

	// Заполняем полученные транзакции
	for i, t := range received {
		history.Received[i] = response.CoinTransaction{
			FromUser: strconv.Itoa(t.FromUserID),
			Amount:   t.Amount,
		}
	}

	return history, nil
}
