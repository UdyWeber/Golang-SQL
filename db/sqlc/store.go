package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Provides all functions to execute db queries with transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// Executes functions in database transactions
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Starts transactions
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})

	// If you could not start transactions return error
	if err != nil {
		return err
	}

	// Gets new instance of Queries Interface
	q := New(tx)

	// Executes the function it is decorating is like a Python decorator
	err = fn(q)

	// If error occurred in Function
	if err != nil {
		// Try to roll back, if it throws an error return a formatted error with rollback error + query error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("[TRANSACTION ERROR]: %v\n[ROLLBACK ERROR]: %v", err, rbErr)
		}
		return err
	}

	// Commit transactions
	return tx.Commit(ctx)
}

// Parameters that are going to be used to perform the operation
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Result of the Transfer operation
type TransferTxResult struct {
	Transfer    Post    `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

var txKey = struct{}{}

// We have to define functions that are going to be run under transaction layer
func (store *Store) TransferTx(ctx context.Context, tp TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Callback function
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "creating transfer")
		// Create transfer between accounts
		result.Transfer, err = q.CreatePost(ctx, CreatePostParams{
			FromAccountID: tp.FromAccountID,
			ToAccountID:   tp.ToAccountID,
			Amount:        tp.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "creating entry to FromAccount")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: tp.FromAccountID,
			Amount:    -tp.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "creating entry to ToAccount")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: tp.ToAccountID,
			Amount:    tp.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "getting to ToAccount")
		account1, err := q.GetAccountForUpdate(ctx, tp.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "updating to ToAccount")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			NewBalance: account1.Balance + tp.Amount,
			ID:         account1.ID,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "getting to FromAccount")
		account2, err := q.GetAccountForUpdate(ctx, tp.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "updating to FromAccount")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			NewBalance: account2.Balance - tp.Amount,
			ID:         account2.ID,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
