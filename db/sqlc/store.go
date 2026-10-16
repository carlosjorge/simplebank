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

		result.Transfer, err = addTransfer(ctx, q, arg)
		if err != nil {
			return err
		}

		result.FromEntry, result.ToEntry, err = addEntries(ctx, q, arg)

		if err != nil {
			return err
		}

		if arg.FromAccount < arg.ToAccount {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccount, -arg.Amount, arg.ToAccount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccount, arg.Amount, arg.FromAccount, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1, account2 Accounts, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	return
}

func addTransfer(
	ctx context.Context,
	q *Queries,
	arg TransferTxParams,
) (transfer Transfers, err error) {
	transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
		FromAccount: arg.FromAccount,
		ToAccount:   arg.ToAccount,
		Amount:      arg.Amount,
	})
	if err != nil {
		return Transfers{}, err
	}
	return transfer, err
}

func addEntries(
	ctx context.Context,
	q *Queries,
	arg TransferTxParams,
) (entry1, entry2 Entries, err error) {
	entry1, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.FromAccount,
		Amount:    -arg.Amount,
	})
	if err != nil {
		return
	}

	entry2, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.ToAccount,
		Amount:    arg.Amount,
	})
	if err != nil {
		return
	}
	return
}
