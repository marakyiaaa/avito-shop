package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"avito_shop/internal/models/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InventoryRepository определяет методы для работы с инвентарем пользователя.
type InventoryRepository interface {
	AddToInventory(ctx context.Context, userID int, itemType string) error
	GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error)
}

// inventoryRepository реализует интерфейс InventoryRepository.
type inventoryRepository struct {
	db *pgxpool.Pool
}

// NewInventoryRepository создает новый экземпляр inventoryRepository.
func NewInventoryRepository(db *pgxpool.Pool) InventoryRepository {
	return &inventoryRepository{db: db}
}

// AddToInventory добавляет товар в инвентарь пользователя.
func (r *inventoryRepository) AddToInventory(ctx context.Context, userID int, itemType string) error {
	const checkQuery = `SELECT quantity FROM inventories WHERE user_id = $1 AND item_type = $2`
	var quantity int
	row := r.db.QueryRow(ctx, checkQuery, userID, itemType)
	err := row.Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			const insertQuery = `INSERT INTO inventories (user_id, item_type, quantity) VALUES ($1, $2, 1)`
			_, err := r.db.Exec(ctx, insertQuery, userID, itemType)
			if err != nil {
				return fmt.Errorf("ошибка при добавлении товара в инвентарь: %w", err)
			}
		} else {
			return fmt.Errorf("ошибка при проверке инвентаря: %w", err)
		}
	} else {
		const updateQuery = `UPDATE inventories SET quantity = $1 WHERE user_id = $2 AND item_type = $3`
		_, err := r.db.Exec(ctx, updateQuery, quantity+1, userID, itemType)
		if err != nil {
			return fmt.Errorf("ошибка при обновлении количества товара: %w", err)
		}
	}
	return nil
}

// GetInventoryByUserID возвращает инвентарь пользователя по его ID.
func (r *inventoryRepository) GetInventoryByUserID(ctx context.Context, userID int) ([]*entities.Inventory, error) {
	const query = `SELECT id, user_id, item_type, quantity FROM inventories WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении инвентаря: %w", err)
	}
	defer rows.Close()

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
