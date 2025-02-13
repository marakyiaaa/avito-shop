package service

import (
	"avito_shop/internal/models/api/response"
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"errors"
	"fmt"
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

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

// SendCoin обновление балансов пользователей
func (s *storeService) SendCoin(ctx context.Context, userID, recipientID int, amount int) error {
	// Получаем отправителя
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Проверяем, достаточно ли средств у отправителя
	if user.Balance < amount {
		return fmt.Errorf("not enough balance")
	}

	// Получаем получателя
	recipient, err := s.userRepo.GetUserByID(ctx, recipientID)
	if err != nil {
		return fmt.Errorf("recipient not found: %w", err)
	}

	// Обновляем баланс отправителя и получателя
	user.Balance -= amount
	recipient.Balance += amount

	// Сохраняем обновленные данные пользователей в базе
	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update sender: %w", err)
	}

	err = s.userRepo.UpdateUser(ctx, recipient)
	if err != nil {
		return fmt.Errorf("failed to update recipient: %w", err)
	}

	// Логика для записи транзакции в таблицу транзакций
	return s.transactionRepo.CreateTransaction(ctx, userID, recipientID, amount)
}
