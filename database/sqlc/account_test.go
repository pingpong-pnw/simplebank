package database

import (
	"context"
	"database/sql"
	"github.com/pingpong-pnw/simplebank/database/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	createParams := CreateAccountParams{
		Owner:    sql.NullString{String: util.RandomOwner(), Valid: true},
		Balance:  sql.NullInt64{Int64: util.RandomMoney(), Valid: true},
		Currency: sql.NullString{String: util.RandomCurrency(), Valid: true},
	}
	account, err := testQuery.CreateAccount(context.Background(), createParams)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, createParams.Owner, account.Owner)
	require.Equal(t, createParams.Balance, account.Balance)
	require.Equal(t, createParams.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestSelectAccountById(t *testing.T) {
	account := createRandomAccount(t)
	selectedAccount, err := testQuery.SelectAccountById(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selectedAccount)
	require.Equal(t, selectedAccount.ID, account.ID)
	require.Equal(t, selectedAccount.Owner, account.Owner)
	require.Equal(t, selectedAccount.Balance, account.Balance)
	require.Equal(t, selectedAccount.Currency, account.Currency)
	require.Equal(t, selectedAccount.CreatedAt, account.CreatedAt)
}

func TestUpdateAccountAccount(t *testing.T) {
	account := createRandomAccount(t)
	updateParams := UpdateAccountParams{
		ID:      account.ID,
		Balance: sql.NullInt64{Int64: util.RandomMoney(), Valid: true},
	}
	updatedAccount, err := testQuery.UpdateAccount(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, updatedAccount.ID, account.ID)
	require.Equal(t, updatedAccount.Owner, account.Owner)
	require.Equal(t, updatedAccount.Balance, updateParams.Balance)
	require.Equal(t, updatedAccount.Currency, account.Currency)
	require.Equal(t, updatedAccount.CreatedAt, account.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	selectedAccount, err := testQuery.SelectAccountById(context.Background(), account.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, selectedAccount)
}

func TestSelectAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	selectAccountParams := SelectAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQuery.SelectAccount(context.Background(), selectAccountParams)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
