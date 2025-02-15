package repository_test

import (
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewUserRepository(db)

	t.Run("успешное создание пользователя", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users \(username, password, balance\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
			WithArgs("user1", "password", 100).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		user := &entities.User{
			Username: "user1",
			Password: "password",
			Balance:  100,
		}

		err := repo.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при создании пользователя", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users \(username, password, balance\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
			WithArgs("user1", "password", 100).
			WillReturnError(errors.New("database error"))

		user := &entities.User{
			Username: "user1",
			Password: "password",
			Balance:  100,
		}

		err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewUserRepository(db)

	t.Run("успешное получение пользователя по ID", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, balance FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance"}).
				AddRow(1, "user1", 100))

		user, err := repo.GetUserByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "user1", user.Username)
		assert.Equal(t, 100, user.Balance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, balance FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при выполнении запроса", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, balance FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewUserRepository(db)

	t.Run("успешное получение пользователя по имени", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password, balance FROM users WHERE username = \$1`).
			WithArgs("user1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "balance"}).
				AddRow(1, "user1", "password", 100))

		user, err := repo.GetUserByUsername(context.Background(), "user1")

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "user1", user.Username)
		assert.Equal(t, "password", user.Password)
		assert.Equal(t, 100, user.Balance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password, balance FROM users WHERE username = \$1`).
			WithArgs("unknown").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUsername(context.Background(), "unknown")

		assert.NoError(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при выполнении запроса", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password, balance FROM users WHERE username = \$1`).
			WithArgs("user1").
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByUsername(context.Background(), "user1")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUserBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мока базы данных: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	repo := repository.NewUserRepository(db)

	t.Run("успешное обновление баланса", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
			WithArgs(200, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUserBalance(context.Background(), 1, 200)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ошибка при обновлении баланса", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
			WithArgs(200, 1).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserBalance(context.Background(), 1, 200)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
