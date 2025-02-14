package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ItemRepository определяет методы для работы с хранилищем элементов.
type ItemRepository interface {
	GetItemByName(ctx context.Context, name string) (*entities.Item, error)
}

// itemRepository реализует интерфейс репозитория для работы с элементами.
// Использует базу данных для хранения и управления данными.
type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

// GetItemByName возвращает элемент по его имени.
func (r *itemRepository) GetItemByName(ctx context.Context, name string) (*entities.Item, error) {
	item := &entities.Item{}
	const query = `SELECT id, name, price FROM items WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)
	err := row.Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}
	return item, nil
}
