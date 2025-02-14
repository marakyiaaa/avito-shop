package repository

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
	"fmt"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, fromUserID int, toUserID int, amount int) (*entities.Transaction, error)
	GetUserTransactions(ctx context.Context, userID int) ([]*entities.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, fromUserID int, toUserID int, amount int) (*entities.Transaction, error) {
	const query = ` INSERT INTO transactions (from_user_id, to_user_id, amount) 
		VALUES ($1, $2, $3) RETURNING id, from_user_id, to_user_id, amount, timestamp`
	tx := &entities.Transaction{}
	err := r.db.QueryRowContext(ctx, query, fromUserID, toUserID, amount).Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}
	return tx, nil
}

func (r *transactionRepository) GetUserTransactions(ctx context.Context, userID int) ([]*entities.Transaction, error) {
	const query = `
		SELECT id, from_user_id, to_user_id, amount, timestamp 
		FROM transactions 
		WHERE from_user_id = $1 OR to_user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		tx := &entities.Transaction{}
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Timestamp); err != nil {
			return nil, fmt.Errorf("error scanning transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
