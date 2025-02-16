package service_test

import (
	"context"
	"fmt"
	"testing"

	"avito_shop/internal/models/entities"
	"avito_shop/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestTransactionService_SendCoins_Success Успешный перевод монет
func TestTransactionService_SendCoins_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	fromUser := &entities.User{ID: 1, Balance: 1000}
	toUser := &entities.User{ID: 2, Balance: 500}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return(toUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 900).Return(nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 2, 600).Return(nil)
	mockTransactionRepo.On("CreateTransaction", mock.Anything, 1, 2, 100).Return(&entities.Transaction{}, nil)

	transactionService := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_SendToSelf Ошибка: перевод самому себе
func TestTransactionService_SendCoins_SendToSelf(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	transactionService := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

	err := transactionService.SendCoins(context.Background(), 1, 1, 100)

	assert.Error(t, err)
	assert.Equal(t, "нельзя отправлять монеты самому себе", err.Error())
}

// TestTransactionService_SendCoins_InsufficientFunds Ошибка: недостаточно средств
func TestTransactionService_SendCoins_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	fromUser := &entities.User{ID: 1, Balance: 50}
	toUser := &entities.User{ID: 2, Balance: 500}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return(toUser, nil)

	transactionService := service.NewTransactionService(mockUserRepo, nil)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "недостаточно средств", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_SenderNotFound Ошибка: отправитель не найден
func TestTransactionService_SendCoins_SenderNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return((*entities.User)(nil), fmt.Errorf("пользователь не найден"))

	transactionService := service.NewTransactionService(mockUserRepo, nil)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "пользователь отправитель не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_ReceiverNotFound Ошибка: получатель не найден
func TestTransactionService_SendCoins_ReceiverNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	fromUser := &entities.User{ID: 1, Balance: 1000}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return((*entities.User)(nil), fmt.Errorf("пользователь не найден"))

	transactionService := service.NewTransactionService(mockUserRepo, nil)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "пользователь получатель не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_UpdateSenderBalanceError Ошибка при обновлении баланса отправителя
func TestTransactionService_SendCoins_UpdateSenderBalanceError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	fromUser := &entities.User{ID: 1, Balance: 1000}
	toUser := &entities.User{ID: 2, Balance: 500}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return(toUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 900).Return(fmt.Errorf("ошибка обновления баланса"))

	transactionService := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "не удалось обновить баланс отправителя: ошибка обновления баланса", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_UpdateReceiverBalanceError Ошибка при обновлении баланса получателя
func TestTransactionService_SendCoins_UpdateReceiverBalanceError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	fromUser := &entities.User{ID: 1, Balance: 1000}
	toUser := &entities.User{ID: 2, Balance: 500}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return(toUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 900).Return(nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 2, 600).Return(fmt.Errorf("ошибка обновления баланса"))

	transactionService := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "не удалось обновить баланс получателя: ошибка обновления баланса", err.Error())
	mockUserRepo.AssertExpectations(t)
}

// TestTransactionService_SendCoins_CreateTransactionError Ошибка при создании транзакции
func TestTransactionService_SendCoins_CreateTransactionError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	fromUser := &entities.User{ID: 1, Balance: 1000}
	toUser := &entities.User{ID: 2, Balance: 500}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(fromUser, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 2).Return(toUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 900).Return(nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 2, 600).Return(nil)
	mockTransactionRepo.On("CreateTransaction", mock.Anything, 1, 2, 100).Return((*entities.Transaction)(nil), fmt.Errorf("ошибка создания транзакции"))

	transactionService := service.NewTransactionService(mockUserRepo, mockTransactionRepo)

	err := transactionService.SendCoins(context.Background(), 1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "не удалось записать транзакцию: ошибка создания транзакции", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}
