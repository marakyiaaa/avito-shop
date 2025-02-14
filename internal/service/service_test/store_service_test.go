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

// TestStoreService_BuyItem_Success Успешная покупка товара.
func TestStoreService_BuyItem_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	expectedItem := &entities.Item{
		ID:    1,
		Name:  "item1",
		Price: 100,
	}
	expectedUser := &entities.User{
		ID:      1,
		Balance: 200,
	}
	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(expectedItem, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 100).Return(nil)
	mockInventoryRepo.On("AddToInventory", mock.Anything, 1, "item1").Return(nil)

	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
	err := storeService.BuyItem(context.Background(), 1, "item1")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

// TestStoreService_BuyItem_ItemNotFound Товар не найден
func TestStoreService_BuyItem_ItemNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(nil, fmt.Errorf("товар не найден"))

	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
	err := storeService.BuyItem(context.Background(), 1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "товар не найден", err.Error())
	mockItemRepo.AssertExpectations(t)
}

//TestStoreService_BuyItem_UserNotFound Пользователь не найден
//func TestStoreService_BuyItem_UserNotFound(t *testing.T) {
//	mockUserRepo := new(MockUserRepository)
//	mockItemRepo := new(MockItemRepository)
//	mockInventoryRepo := new(MockInventoryRepository)
//
//	expectedItem := &entities.Item{
//		ID:    1,
//		Name:  "item1",
//		Price: 100,
//	}
//	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(expectedItem, nil)
//	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(nil, fmt.Errorf("пользователь не найден"))
//
//	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
//	err := storeService.BuyItem(context.Background(), 1, "item1")
//
//	assert.Error(t, err)
//	assert.Equal(t, "пользователь не найден", err.Error())
//	mockItemRepo.AssertExpectations(t)
//	mockUserRepo.AssertExpectations(t)
//}

// TestStoreService_BuyItem_InsufficientFunds Недостаточно средств
func TestStoreService_BuyItem_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	expectedItem := &entities.Item{
		ID:    1,
		Name:  "item1",
		Price: 100,
	}
	expectedUser := &entities.User{
		ID:      1,
		Balance: 50,
	}
	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(expectedItem, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)

	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
	err := storeService.BuyItem(context.Background(), 1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "недостаточно средств", err.Error())
	mockItemRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestStoreService_BuyItem_UpdateBalanceError Ошибка при обновлении баланса
func TestStoreService_BuyItem_UpdateBalanceError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	expectedItem := &entities.Item{
		ID:    1,
		Name:  "item1",
		Price: 100,
	}
	expectedUser := &entities.User{
		ID:      1,
		Balance: 200,
	}
	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(expectedItem, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 100).Return(fmt.Errorf("ошибка при обновлении баланса"))

	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
	err := storeService.BuyItem(context.Background(), 1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "не удалось обновить баланс: ошибка при обновлении баланса", err.Error())
	mockItemRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestStoreService_BuyItem_AddToInventoryError Ошибка при добавлении в инвентарь
func TestStoreService_BuyItem_AddToInventoryError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	// Настраиваем моки
	expectedItem := &entities.Item{
		ID:    1,
		Name:  "item1",
		Price: 100,
	}
	expectedUser := &entities.User{
		ID:      1,
		Balance: 200,
	}
	mockItemRepo.On("GetItemByName", mock.Anything, "item1").Return(expectedItem, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, 1, 100).Return(nil)
	mockInventoryRepo.On("AddToInventory", mock.Anything, 1, "item1").Return(fmt.Errorf("ошибка при добавлении в инвентарь"))

	storeService := service.NewStoreService(mockUserRepo, mockItemRepo, mockInventoryRepo)
	err := storeService.BuyItem(context.Background(), 1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "не удалось добавить предмет в инвентарь: ошибка при добавлении в инвентарь", err.Error())
	mockItemRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}
