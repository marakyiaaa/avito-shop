package service_test

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/service"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// TestInfoService_GetUserInfo_Success Успешное получение информации о пользователе
func TestInfoService_GetUserInfo_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	expectedUser := &entities.User{ID: 1, Username: "user1", Balance: 1000}
	expectedInventory := []*entities.Inventory{{ID: 1, UserID: 1, ItemType: "type1", Quantity: 1}}
	expectedTransactions := []*entities.Transaction{{FromUserID: 1, ToUserID: 2, Amount: 100}}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
	mockInventoryRepo.On("GetInventoryByUserID", mock.Anything, 1).Return(expectedInventory, nil)
	mockTransactionRepo.On("GetUserTransactions", mock.Anything, 1).Return(expectedTransactions, nil)

	infoService := service.NewInfoService(mockUserRepo, mockItemRepo, mockInventoryRepo, mockTransactionRepo)

	response, err := infoService.GetUserInfo(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedUser, response.User)
	assert.Equal(t, expectedInventory, response.Inventory)
	assert.Equal(t, expectedTransactions, response.Transactions)

	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

// TestInfoService_GetUserInfo_UserNotFound Ошибка при получении пользователя
func TestInfoService_GetUserInfo_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return((*entities.User)(nil), fmt.Errorf("пользователь не найден"))

	infoService := service.NewInfoService(mockUserRepo, mockItemRepo, mockInventoryRepo, mockTransactionRepo)

	response, err := infoService.GetUserInfo(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockUserRepo.AssertExpectations(t)
}

// TestInfoService_GetUserInfo_InventoryError Ошибка при получении инвентаря
func TestInfoService_GetUserInfo_InventoryError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	expectedUser := &entities.User{ID: 1, Username: "user1", Balance: 1000}

	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
	mockInventoryRepo.On("GetInventoryByUserID", mock.Anything, 1).Return([]*entities.Inventory{}, fmt.Errorf("ошибка инвентаря"))

	infoService := service.NewInfoService(mockUserRepo, mockItemRepo, mockInventoryRepo, mockTransactionRepo)

	response, err := infoService.GetUserInfo(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

// TestInfoService_GetUserInfo_TransactionsError Ошибка при получении транзакций
//func TestInfoService_GetUserInfo_TransactionsError(t *testing.T) {
//	mockUserRepo := new(MockUserRepository)
//	mockItemRepo := new(MockItemRepository)
//	mockInventoryRepo := new(MockInventoryRepository)
//	mockTransactionRepo := new(MockTransactionRepository)
//
//	// Ожидаемые данные
//	expectedUser := &entities.User{ID: 1, Username: "user1", Balance: 1000}
//	expectedInventory := []*entities.Inventory{{ID: 1, UserID: 1, ItemType: "type1", Quantity: 1}}
//
//	// Настройка моков
//	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
//	mockInventoryRepo.On("GetInventoryByUserID", mock.Anything, 1).Return(expectedInventory, nil)
//	mockTransactionRepo.On("GetUserTransactions", mock.Anything, 1).Return([]*entities.Transaction{}, fmt.Errorf("ошибка транзакций"))
//
//	// Создаем сервис
//	infoService := service.NewInfoService(mockUserRepo, mockItemRepo, mockInventoryRepo, mockTransactionRepo)
//
//	// Вызываем метод
//	response, err := infoService.GetUserInfo(context.Background(), 1)
//
//	// Проверяем результат
//	assert.Error(t, err)
//	assert.Nil(t, response)
//	mockUserRepo.AssertExpectations(t)
//	mockInventoryRepo.AssertExpectations(t)
//	mockTransactionRepo.AssertExpectations(t)
//}
