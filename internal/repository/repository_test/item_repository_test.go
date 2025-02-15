package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"

	"avito_shop/internal/repository"
)

func TestGetItemByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewItemRepository(db)

	t.Run("успешное получение элемента", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name, price FROM items WHERE name = \$1`).
			WithArgs("cup").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow(1, "cup", 100))

		item, err := repo.GetItemByName(context.Background(), "cup")

		assert.NoError(t, err)
		assert.Equal(t, 1, item.ID)
		assert.Equal(t, "cup", item.Name)
		assert.Equal(t, 100, item.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("элемент не найден", func(t *testing.T) {
		// Ожидаем запрос на получение элемента
		mock.ExpectQuery(`SELECT id, name, price FROM items WHERE name = \$1`).
			WithArgs("unknown").
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetItemByName(context.Background(), "unknown")

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "item not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при выполнении запроса", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name, price FROM items WHERE name = \$1`).
			WithArgs("cup").
			WillReturnError(errors.New("database error"))

		item, err := repo.GetItemByName(context.Background(), "cup")

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
