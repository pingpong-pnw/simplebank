package database

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before:", account1.Balance.Int64, account2.Balance.Int64)
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.CreateTransferTx(context.Background(), TransferTxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existing := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testQuery.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, (-1)*amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = testQuery.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = testQuery.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">> Transaction:", fromAccount.Balance.Int64, toAccount.Balance.Int64)
		diffBalance1 := account1.Balance.Int64 - fromAccount.Balance.Int64
		diffBalance2 := toAccount.Balance.Int64 - account2.Balance.Int64
		require.Equal(t, diffBalance1, diffBalance2)
		require.True(t, diffBalance1 > 0)
		require.True(t, diffBalance2%amount == 0)

		k := int(diffBalance1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existing, k)
		existing[k] = true
	}
	updateAccount1, err := testQuery.SelectAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQuery.SelectAccountById(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> After:", updateAccount1.Balance.Int64, updateAccount2.Balance.Int64)
	require.Equal(t, account1.Balance.Int64-int64(n)*amount, updateAccount1.Balance.Int64)
	require.Equal(t, account2.Balance.Int64+int64(n)*amount, updateAccount2.Balance.Int64)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before:", account1.Balance.Int64, account2.Balance.Int64)
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.CreateTransferTx(context.Background(), TransferTxRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	updateAccount1, err := testQuery.SelectAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQuery.SelectAccountById(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> After:", updateAccount1.Balance.Int64, updateAccount2.Balance.Int64)
	require.Equal(t, account1.Balance.Int64, updateAccount1.Balance.Int64)
	require.Equal(t, account2.Balance.Int64, updateAccount2.Balance.Int64)
}
