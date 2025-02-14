package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type ItemRepository interface {
	GetItemByName(ctx context.Context, name string) (*entities.Item, error)
	AddToInventory(ctx context.Context, userID int, itemType string) error
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

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

func (r *itemRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
	const checkQuery = `SELECT quantity FROM inventories WHERE user_id = $1 AND item_type = $2`
	var quantity int
	row := r.db.QueryRowContext(ctx, checkQuery, userID, itemType)
	err := row.Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			const insertQuery = `INSERT INTO inventories (user_id, item_type, quantity) VALUES ($1, $2, 1)`
			_, err := r.db.ExecContext(ctx, insertQuery, userID, itemType)
			if err != nil {
				return fmt.Errorf("ошибка при добавлении товара в инвентарь: %w", err)
			}
		} else {
			return fmt.Errorf("ошибка при проверке инвентаря: %w", err)
		}
	} else {
		const updateQuery = `UPDATE inventories SET quantity = $1 WHERE user_id = $2 AND item_type = $3`
		_, err := r.db.ExecContext(ctx, updateQuery, quantity+1, userID, itemType)
		if err != nil {
			return fmt.Errorf("ошибка при обновлении количества товара: %w", err)
		}
	}
	return nil
}
