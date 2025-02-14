package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// InventoryRepository определяет методы для работы с инвентарем пользователя.
type InventoryRepository interface {
	AddToInventory(ctx context.Context, userID int, itemType string) error
	GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error)
}

// inventoryRepository реализует интерфейс InventoryRepository.
type inventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository создает новый экземпляр inventoryRepository.
func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// AddToInventory добавляет элемент в инвентарь пользователя.
func (r *inventoryRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
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

// GetInventoryByUserID возвращает инвентарь пользователя по его ID.
func (r *inventoryRepository) GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error) {
	const query = `SELECT id, user_id, item_type, quantity FROM inventories WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении инвентаря: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Логирование ошибки закрытия rows (опционально)
		}
	}(rows)

	var inventories []*entities.Inventory
	for rows.Next() {
		inventory := &entities.Inventory{}
		if err := rows.Scan(&inventory.ID, &inventory.UserID, &inventory.ItemType, &inventory.Quantity); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании инвентаря: %w", err)
		}
		inventories = append(inventories, inventory)
	}
	return inventories, nil
}
