package database

import (
	"context"
	"github.com/pingpong-pnw/simplebank/database/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	amount := util.RandomAmount()
	entry, err := testQuery.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	})
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)
	require.True(t, entry.CreatedAt.Valid)
	return entry
}

func createFixEntry(t *testing.T, account Account) Entry {
	amount := util.RandomAmount()
	entry, err := testQuery.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	})
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)
	require.True(t, entry.CreatedAt.Valid)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)
	selectedEntry, err := testQuery.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selectedEntry)
	require.Equal(t, entry.ID, selectedEntry.ID)
	require.Equal(t, entry.AccountID, selectedEntry.AccountID)
	require.Equal(t, entry.Amount, selectedEntry.Amount)
	require.Equal(t, entry.CreatedAt, selectedEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	n := 10
	for i := 0; i < n; i++ {
		createFixEntry(t, account)
	}

	entries, err := testQuery.ListEntries(context.Background(), ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	})
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, account.ID, entry.AccountID)
		require.NotZero(t, entry.Amount)
		require.NotZero(t, entry.CreatedAt)
	}
}
