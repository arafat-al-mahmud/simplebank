package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(db)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("-> before : ", account1.Balance, account2.Balance)

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
	existed := make(map[int]bool)
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

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check account's balance
		fmt.Println(" tx : ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%arg.Amount == 0)

		k := diff1 / arg.Amount
		require.True(t, k >= 1 && k <= int64(n))
		require.NotContains(t, existed, k)
		existed[int(k)] = true
	}

	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("-> after : ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*arg.Amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*arg.Amount, updatedAccount2.Balance)
}

func TestTransferTxdDeadlock(t *testing.T) {
	store := NewStore(db)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("account1 : ", account1.ID, " account2 : ", account2.ID)
	fmt.Println("-> before : ", account1.Balance, account2.Balance)

	//run n concurrent transfer transactions
	n := 10
	transferAmount := 1000
	errorChannel := make(chan error)

	for i := 0; i < n; i++ {

		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		arg := TransferTxParams{
			FromAccountID: fromAccountId,
			ToAccountID:   toAccountId,
			Amount:        int64(transferAmount),
		}

		go func() {
			_, err := store.TransferTx(context.Background(), arg)

			errorChannel <- err

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorChannel
		require.NoError(t, err)
	}

	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("-> after : ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
