package service_test

import (
	"avito_shop/internal/models/entities"
	"context"
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
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.User), args.Error(1)
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
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
	args := m.Called(ctx, userID, itemType)
	return args.Error(0)
}

func (m *MockItemRepository) GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Inventory), args.Error(1)
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
	return args.Get(0).([]*entities.Inventory), args.Error(1)
}
