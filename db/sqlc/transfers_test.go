package db

import (
	"context"
	"testing"

	"github.com/carlosjorge/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfers {
	account1_args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account1, err := testQueries.CreateAccount(context.Background(), account1_args)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	account2_args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account2, err := testQueries.CreateAccount(context.Background(), account2_args)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	args := CreateTransferParams{
		FromAccount: account1.ID,
		ToAccount:   account2.ID,
		Amount:      util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	require.NotEmpty(t, transfer)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:          transfer.ID,
		FromAccount: transfer.FromAccount,
		ToAccount:   transfer.ToAccount,
		Amount:      util.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, arg.Amount, transfer2.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for range 10 {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
