package database

import (
	"context"
	"github.com/pingpong-pnw/simplebank/database/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(t *testing.T, transferParams CreateTransferParams) Transfer {
	transfer, err := testQuery.CreateTransfer(context.Background(), transferParams)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.Equal(t, transferParams.FromAccountID, transfer.FromAccountID)
	require.Equal(t, transferParams.ToAccountID, transferParams.ToAccountID)
	require.Equal(t, transferParams.Amount, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)
	require.True(t, transfer.CreatedAt.Valid)
	return transfer
}

func createRandomTransferParam(t *testing.T) CreateTransferParams {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	return CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}
}

func createFixTransferParams(t *testing.T, account1 Account, account2 Account) CreateTransferParams {
	return CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, createRandomTransferParam(t))
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t, createRandomTransferParam(t))
	selectedTransfer, err := testQuery.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selectedTransfer)
	require.Equal(t, transfer.ID, selectedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, selectedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, selectedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, selectedTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, selectedTransfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 10
	for i := 0; i < n; i++ {
		createRandomTransfer(t, createFixTransferParams(t, account1, account2))
	}
	listTransferParams := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQuery.ListTransfers(context.Background(), listTransferParams)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfers)
		require.NotZero(t, transfer.ID)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.NotZero(t, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
	}
}
