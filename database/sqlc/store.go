package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) CreateTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) CreateTransferTx(ctx context.Context, request TransferTxRequest) (TransferTxResult, error) {
	var transferTxResult TransferTxResult
	err := s.CreateTransaction(ctx, func(queries *Queries) error {
		var err error
		transferTxResult.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: request.FromAccountID,
			ToAccountID:   request.ToAccountID,
			Amount:        request.Amount,
		})
		if err != nil {
			return nil
		}
		transferTxResult.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: request.FromAccountID,
			Amount:    (-1) * request.Amount,
		})
		if err != nil {
			return err
		}
		transferTxResult.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: request.ToAccountID,
			Amount:    request.Amount,
		})
		if err != nil {
			return err
		}
		if request.FromAccountID < request.ToAccountID {
			transferTxResult.FromAccount, transferTxResult.ToAccount, err = UpdateAccountBalance(ctx, queries, request.FromAccountID, (-1)*request.Amount, request.ToAccountID, request.Amount)
		} else {
			transferTxResult.ToAccount, transferTxResult.FromAccount, err = UpdateAccountBalance(ctx, queries, request.ToAccountID, request.Amount, request.FromAccountID, (-1)*request.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return transferTxResult, err
}

func UpdateAccountBalance(ctx context.Context, queries *Queries, fromAccountID int64, fromAccountAmount int64, toAccountID int64, toAccountAmount int64) (account1 Account, account2 Account, err error) {
	account1, err = queries.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: fromAccountID,
		Amount: sql.NullInt64{
			Int64: fromAccountAmount,
			Valid: true,
		},
	})
	if err != nil {
		return
	}
	account2, err = queries.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: toAccountID,
		Amount: sql.NullInt64{
			Int64: toAccountAmount,
			Valid: true,
		},
	})
	if err != nil {
		return
	}
	return
}
