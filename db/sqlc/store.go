package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Stores provide all functions to execute queries and transactions.
type Store struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTx executes a function within a database transaction.
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccount int64 `json:"from_account_id"`
	ToAccount   int64 `json:"to_account_id"`
	Amount      int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// TransferTX performs a money tranfer from one account to another.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction.
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccount: arg.FromAccount,
			ToAccount:   arg.ToAccount,
			Amount:      arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccount,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccount,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// if arg.FromAccount < arg.ToAccount {
		// 	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.FromAccount,
		// 		Amount: -arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.ToAccount,
		// 		Amount: arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.ToAccount,
		// 		Amount: arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.FromAccount,
		// 		Amount: -arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})
	return result, err
}
