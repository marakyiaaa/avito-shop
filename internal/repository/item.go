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
	AddToInventory(ctx context.Context, userID, itemID int) error
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

//// AddToInventory запись о покупке в таблицу user_inventory
//func (r *itemRepository) AddToInventory(ctx context.Context, userID, itemID int) error {
//	const query = `INSERT INTO user_inventory (user_id, item_id) VALUES ($1, $2)`
//	_, err := r.db.ExecContext(ctx, query, userID, itemID)
//	if err != nil {
//		return fmt.Errorf("ошибка при добавлении предмета в инвентарь: %w", err)
//	}
//	return nil
//}

// AddToInventory запись о покупке в таблицу user_inventory
func (r *itemRepository) AddToInventory(ctx context.Context, userID, itemID int) error {
	// Проверяем, есть ли уже товар в инвентаре
	const checkQuery = `SELECT quantity FROM user_inventory WHERE user_id = $1 AND item_id = $2`
	var quantity int
	row := r.db.QueryRowContext(ctx, checkQuery, userID, itemID)
	err := row.Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Товара нет, добавляем новый
			const insertQuery = `INSERT INTO user_inventory (user_id, item_id, quantity) VALUES ($1, $2, $3)`
			_, err := r.db.ExecContext(ctx, insertQuery, userID, itemID)
			if err != nil {
				return fmt.Errorf("ошибка при добавлении предмета в инвентарь: %w", err)
			}
		} else {
			return fmt.Errorf("ошибка при проверке товара в инвентаре: %w", err)
		}
	} else {
		// Товар есть, увеличиваем количество
		const updateQuery = `UPDATE user_inventory SET quantity = $1 WHERE user_id = $2 AND item_id = $3`
		_, err := r.db.ExecContext(ctx, updateQuery, quantity+1, userID, itemID)
		if err != nil {
			return fmt.Errorf("ошибка при обновлении количества товара: %w", err)
		}
	}
	return nil
}
