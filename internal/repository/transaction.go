package repository

//хз новое

import (
	"avito_shop/internal/models/entities"
	"context"
	"database/sql"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, userID, itemID int, price int) (*entities.Transaction, error)
	GetSentTransactions(ctx context.Context, userID int) ([]entities.Transaction, error)
	GetReceivedTransactions(ctx context.Context, userID int) ([]entities.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, fromUserID, toUserID, amount int) (*entities.Transaction, error) {
	const query = `INSERT INTO transactions (from_user_id, to_user_id, amount) 
                   VALUES ($1, $2, $3) RETURNING id, from_user_id, to_user_id, amount, timestamp`
	transaction := &entities.Transaction{}
	err := r.db.QueryRowContext(ctx, query, fromUserID, toUserID, amount).
		Scan(&transaction.ID, &transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.Timestamp)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

//func (r *transactionRepository) GetSentTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
//	var transactions []entities.Transaction
//	query := `SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE from_user_id = $1`
//	err := r.db.SelectContext(ctx, &transactions, query, userID)
//	return transactions, err
//}
//
//func (r *transactionRepository) GetReceivedTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
//	var transactions []entities.Transaction
//	query := `SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE to_user_id = $1`
//	err := r.db.SelectContext(ctx, &transactions, query, userID)
//	return transactions, err
//}

func (r *transactionRepository) GetSentTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	query := `SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE from_user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetReceivedTransactions(ctx context.Context, userID int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	query := `SELECT id, from_user_id, to_user_id, amount, timestamp FROM transactions WHERE to_user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
