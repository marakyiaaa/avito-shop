package service_test

import (
	"context"
	"errors"

	"avito_shop/internal/models/entities"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository — мок для UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	args := m.Called(ctx, id)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	user, ok := args.Get(0).(*entities.User)
	if !ok {
		return nil, errors.New("failed to cast to *entities.User")
	}
	return user, nil
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	user, ok := args.Get(0).(*entities.User)
	if !ok {
		return nil, errors.New("failed to cast to *entities.User")
	}
	return user, nil
}

func (m *MockUserRepository) UpdateUserBalance(ctx context.Context, id int, newBalance int) error {
	args := m.Called(ctx, id, newBalance)
	return args.Error(0)
}

// MockItemRepository — мок для ItemRepository.
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetItemByName(ctx context.Context, name string) (*entities.Item, error) {
	args := m.Called(ctx, name)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	item, ok := args.Get(0).(*entities.Item)
	if !ok {
		return nil, errors.New("failed to cast to *entities.Item")
	}
	return item, nil
}

func (m *MockItemRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
	args := m.Called(ctx, userID, itemType)
	return args.Error(0)
}

func (m *MockItemRepository) GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error) {
	args := m.Called(ctx, userID)

	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	inventory, ok := args.Get(0).([]*entities.Inventory)
	if !ok {
		return nil, errors.New("failed to cast to []*entities.Inventory")
	}
	return inventory, nil
}

// MockTransactionRepository — мок для TransactionRepository.
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetUserTransactions(ctx context.Context, userID int) ([]*entities.Transaction, error) {
	args := m.Called(ctx, userID)

	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	transactions, ok := args.Get(0).([]*entities.Transaction)
	if !ok {
		return nil, errors.New("failed to cast to []*entities.Transaction")
	}
	return transactions, nil
}
func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, fromUserID int, toUserID int, amount int) (*entities.Transaction, error) {
	args := m.Called(ctx, fromUserID, toUserID, amount)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	transaction, ok := args.Get(0).(*entities.Transaction)
	if !ok {
		return nil, errors.New("failed to cast to *entities.Transaction")
	}
	return transaction, nil
}

// MockInventoryRepository — мок для InventoryRepository.
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
	args := m.Called(ctx, userID, itemType)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error) {
	args := m.Called(ctx, userID)

	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	inventory, ok := args.Get(0).([]*entities.Inventory)
	if !ok {
		return nil, errors.New("failed to cast to []*entities.Inventory")
	}
	return inventory, nil
}
