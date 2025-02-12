package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type ItemRepository interface {
	GetItemByName(ctx context.Context, name string) (*entities.Item, error)
	GetAllItems(ctx context.Context) ([]entities.Item, error)
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

// GetItemByName находит предмет по названию
func (r *itemRepository) GetItemByName(ctx context.Context, name string) (*entities.Item, error) {
	item := &entities.Item{}
	const query = `SELECT id, name, price FROM items WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)
	err := row.Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return item, nil
}

// GetAllItems получает список всех доступных предметов
func (r *itemRepository) GetAllItems(ctx context.Context) ([]entities.Item, error) {

	const query = `SELECT id, name, price FROM items`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
