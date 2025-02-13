package service

import (
	"avito_shop/internal/repository"
	"fmt"
)

type TransactionService struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// Отправить монеты от одного пользователя к другому
func (s *TransactionService) SendCoins(senderID, recipientID int, amount float64) error {
	// Получаем отправителя
	sender, err := s.userRepo.GetUserByID(senderID)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}

	// Проверяем баланс отправителя
	if sender.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// Получаем получателя
	recipient, err := s.userRepo.GetUserByID(recipientID)
	if err != nil {
		return fmt.Errorf("recipient not found: %w", err)
	}

	// Обновляем баланс отправителя и получателя
	sender.Balance -= amount
	recipient.Balance += amount

	// Сохраняем обновления в БД
	err = s.userRepo.UpdateUser(sender)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}
	err = s.userRepo.UpdateUser(recipient)
	if err != nil {
		return fmt.Errorf("failed to update recipient balance: %w", err)
	}

	// Записываем транзакцию
	transaction := &models.Transaction{
		SenderID:    senderID,
		RecipientID: recipientID,
		Amount:      amount,
	}
	err = s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
