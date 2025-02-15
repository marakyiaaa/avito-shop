package repository_test

import (
	"avito_shop/internal/repository"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddToInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewInventoryRepository(db)

	t.Run("добавление нового товара", func(t *testing.T) {
		mock.ExpectQuery(`SELECT quantity FROM inventories WHERE user_id = \$1 AND item_type = \$2`).
			WithArgs(1, "cup").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec(`INSERT INTO inventories \(user_id, item_type, quantity\) VALUES \(\$1, \$2, 1\)`).
			WithArgs(1, "cup").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddToInventory(context.Background(), 1, "cup")

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("обновление количества существующего товара", func(t *testing.T) {
		mock.ExpectQuery(`SELECT quantity FROM inventories WHERE user_id = \$1 AND item_type = \$2`).
			WithArgs(1, "cup").
			WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(5))

		mock.ExpectExec(`UPDATE inventories SET quantity = \$1 WHERE user_id = \$2 AND item_type = \$3`).
			WithArgs(6, 1, "cup").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddToInventory(context.Background(), 1, "cup")

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetInventoryByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewInventoryRepository(db)

	t.Run("пустой инвентарь", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, item_type, quantity FROM inventories WHERE user_id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "item_type", "quantity"}))

		result, err := repo.GetInventoryByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("инвентарь с несколькими товарами", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, item_type, quantity FROM inventories WHERE user_id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "item_type", "quantity"}).
				AddRow(1, 1, "cup", 5).
				AddRow(2, 1, "pen", 10))

		result, err := repo.GetInventoryByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "cup", result[0].ItemType)
		assert.Equal(t, 5, result[0].Quantity)
		assert.Equal(t, "pen", result[1].ItemType)
		assert.Equal(t, 10, result[1].Quantity)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
