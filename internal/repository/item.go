package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"avito_shop/internal/models/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ItemRepository определяет методы для работы с хранилищем элементов.
type ItemRepository interface {
	GetItemByName(ctx context.Context, name string) (*entities.Item, error)
}

// itemRepository реализует интерфейс репозитория для работы с элементами.
// Использует базу данных для хранения и управления данными.
type itemRepository struct {
	db *pgxpool.Pool
}

// NewItemRepository создает новый экземпляр itemRepository.
func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db}
}

// GetItemByName возвращает товар по его имени.
func (r *itemRepository) GetItemByName(ctx context.Context, name string) (*entities.Item, error) {
	item := &entities.Item{}
	const query = `SELECT id, name, price FROM items WHERE name = $1`
	row := r.db.QueryRow(ctx, query, name)
	err := row.Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("товар не найден")
		}
		return nil, err
	}
	return item, nil
}
