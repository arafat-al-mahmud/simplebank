package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(db)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        1000,
	}

	//run n concurrent transfer transactions
	n := 5
	resultChannel := make(chan TransferTxResult)
	errorChannel := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), arg)

			errorChannel <- err
			resultChannel <- result

		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errorChannel
		require.NoError(t, err)

		result := <-resultChannel
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, arg.Amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		fromAccountEntry := result.FromEntry
		require.NotEmpty(t, fromAccountEntry)
		require.Equal(t, account1.ID, fromAccountEntry.AccountID)
		require.Equal(t, -arg.Amount, fromAccountEntry.Amount)
		require.NotZero(t, fromAccountEntry.ID)
		require.NotZero(t, fromAccountEntry.CreatedAt)

		toAccountEntry := result.ToEntry
		require.NotEmpty(t, toAccountEntry)
		require.Equal(t, account2.ID, toAccountEntry.AccountID)
		require.Equal(t, arg.Amount, toAccountEntry.Amount)
		require.NotZero(t, toAccountEntry.ID)
		require.NotZero(t, toAccountEntry.CreatedAt)
	}
}
