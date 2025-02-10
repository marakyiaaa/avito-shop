package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (usename, password, coins)
              VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Coins, user.Password, time.Now(), time.Now()).Scan(&user.ID)
	return err
}

//func (r *UserRepository) CreateUser(user *entities.User) error {
//	query := `INSERT INTO users (username, password, coins) VALUES ($1, $2, $3) RETURNING id`
//	return r.db.QueryRow(query, user.Username, user.Password, user.Coins).Scan(&user.ID)
//}
//
//func (r *UserRepository) GetUserByID(id int) (*entities.User, error) {
//	user := &entities.User{}
//	query := `SELECT id, username, coins FROM users WHERE id = $1`
//	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Coins)
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}
