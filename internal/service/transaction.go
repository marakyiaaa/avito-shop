package service

import (
	"avito_shop/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
)

type TransactionService interface {
	SendCoins(ctx context.Context, fromUserID int, toUserID int, amount int) error
}

type transactionService struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: transactionRepo}
}

// SendCoins осуществляет отправку монет другому пользователю
func (s *transactionService) SendCoins(ctx context.Context, fromUserID int, toUserID int, amount int) error {
	fromUser, err := s.userRepo.GetUserByID(ctx, fromUserID)
	if err != nil || fromUser == nil {
		log.Println("Ошибка при получении отправителя")
		return errors.New("пользователь отправитель не найден")
	}

	toUser, err := s.userRepo.GetUserByID(ctx, toUserID)
	if err != nil || toUser == nil {
		log.Println("Ошибка при получении получателя")
		return errors.New("пользователь получатель не найден")
	}

	if fromUser.Balance < amount {
		log.Println("Недостаточно средств")
		return errors.New("недостаточно средств")
	}

	newBalance := fromUser.Balance - amount
	err = s.userRepo.UpdateUserBalance(ctx, fromUserID, newBalance)
	if err != nil {
		log.Println("Ошибка при обновлении баланса отправителя:", err)
		return fmt.Errorf("не удалось обновить баланс отправителя: %w", err)
	}

	err = s.userRepo.UpdateUserBalance(ctx, toUserID, toUser.Balance+amount)
	if err != nil {
		log.Println("Ошибка при обновлении баланса получателя:", err)
		return fmt.Errorf("не удалось обновить баланс получателя: %w", err)
	}

	_, err = s.transactionRepo.CreateTransaction(ctx, fromUserID, toUserID, amount)
	if err != nil {
		log.Println("Ошибка при записи транзакции:", err)
		return fmt.Errorf("не удалось записать транзакцию: %w", err)
	}

	log.Println("Перевод завершен успешно")
	return nil
}
