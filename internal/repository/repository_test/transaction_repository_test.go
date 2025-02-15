package repository_test

import (
	"avito_shop/internal/repository"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewTransactionRepository(db)

	t.Run("успешное создание транзакции", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO transactions \(from_user_id, to_user_id, amount\) VALUES \(\$1, \$2, \$3\) RETURNING id, from_user_id, to_user_id, amount, timestamp`).
			WithArgs(1, 2, 100).
			WillReturnRows(sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "timestamp"}).
				AddRow(1, 1, 2, 100, time.Now()))

		tx, err := repo.CreateTransaction(context.Background(), 1, 2, 100)

		assert.NoError(t, err)
		assert.Equal(t, 1, tx.ID)
		assert.Equal(t, 1, tx.FromUserID)
		assert.Equal(t, 2, tx.ToUserID)
		assert.Equal(t, 100, tx.Amount)
		assert.NotNil(t, tx.Timestamp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при создании транзакции", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO transactions \(from_user_id, to_user_id, amount\) VALUES \(\$1, \$2, \$3\) RETURNING id, from_user_id, to_user_id, amount, timestamp`).
			WithArgs(1, 2, 100).
			WillReturnError(sql.ErrConnDone)

		tx, err := repo.CreateTransaction(context.Background(), 1, 2, 100)

		assert.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "error creating transaction")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewTransactionRepository(db)

	t.Run("успешное получение списка транзакций", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE from_user_id = \$1 OR to_user_id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "timestamp"}).
				AddRow(1, 1, 2, 100, time.Now()).
				AddRow(2, 3, 1, 50, time.Now()))

		transactions, err := repo.GetUserTransactions(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, transactions, 2)
		assert.Equal(t, 1, transactions[0].ID)
		assert.Equal(t, 1, transactions[0].FromUserID)
		assert.Equal(t, 2, transactions[0].ToUserID)
		assert.Equal(t, 100, transactions[0].Amount)
		assert.NotNil(t, transactions[0].Timestamp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("пустой список транзакций", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE from_user_id = \$1 OR to_user_id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "timestamp"}))

		transactions, err := repo.GetUserTransactions(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, transactions)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при выполнении запроса", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE from_user_id = \$1 OR to_user_id = \$1`).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		transactions, err := repo.GetUserTransactions(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.Contains(t, err.Error(), "error fetching transactions")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
