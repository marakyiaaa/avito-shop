package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	UpdateUserBalance(ctx context.Context, id int, newBalance int) error
	//GetAllUsers(ctx context.Context) ([]entities.User, error)
}

// Структура, реализующая интерфейс
type userRepository struct {
	db *sql.DB
}

// NewUserRepository Лучше возвращать интерфейс UserRepository, а не *userRepository - это даёт больше гибкости
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser Регистрация пользователя
func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	const query = `INSERT INTO users (username, password, coins)
              VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Coins).Scan(&user.ID)
	return err
}

// GetUserByID Получить пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, username, coins FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // если пользователя нет, возвращаем nil
	}
	if err != nil {
		return nil, err //ошибки попозже
	}
	return user, nil
}

// GetUserByUsername Получить пользователя по логину
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user := &entities.User{}
	const query = `SELECT id, username, password, coins FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // если пользователя нет, возвращаем nil
	}
	if err != nil {
		log.Println("Ошибка при запросе")
		return nil, err //ошибки попозже
	}
	return user, nil
}

// UpdateUserBalance Обновить баланс пользователя
func (r *userRepository) UpdateUserBalance(ctx context.Context, id int, newBalance int) error {
	query := `UPDATE users SET coins = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newBalance, id)
	return err
}

//получить всех пользователей
//func (r *userRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
//	query := `SELECT id, username, coins FROM users`
//	rows, err := r.db.QueryContext(ctx, query)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var users []entities.User
//	for rows.Next() {
//		var user entities.User
//		if err := rows.Scan(&user.ID, &user.Username, &user.Coins); err != nil {
//			return nil, err
//		}
//		users = append(users, user)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return users, nil
//}

//// Регистрация пользователя с использованием Squirrel
//func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
//	query, args, err := squirrel.Insert("users").
//		Columns("username", "password", "coins").
//		Values(user.Username, user.Password, user.Coins).
//		Suffix("RETURNING id").
//		PlaceholderFormat(squirrel.Dollar).
//		ToSql()
//	if err != nil {
//		return err
//	}
//	return r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
//}
