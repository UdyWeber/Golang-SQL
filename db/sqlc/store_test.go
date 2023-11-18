package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfertTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	args := TransferTxParams{
		FromAccountID: account2.ID,
		ToAccountID:   account1.ID,
		Amount:        amount,
	}

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("transaction %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := testStore.TransferTx(ctx, args)

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.ToAccount)

		require.Equal(t, result.ToAccount.ID, account1.ID)
		require.Equal(t, result.ToAccount.Owner, account1.Owner)
		require.Equal(t, result.ToAccount.Currency, account1.Currency)

		require.NotEmpty(t, result.FromAccount)

		require.Equal(t, result.FromAccount.ID, account2.ID)
		require.Equal(t, result.FromAccount.Owner, account2.Owner)
		require.Equal(t, result.FromAccount.Currency, account2.Currency)

		require.NotEmpty(t, result.ToEntry)

		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)
		require.Equal(t, result.ToEntry.Amount, args.Amount)

		require.NotEmpty(t, result.FromEntry)

		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.Equal(t, result.FromEntry.Amount, -args.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		require.Equal(t, transfer.Amount, args.Amount)
		require.Equal(t, transfer.ToAccountID, args.ToAccountID)
		require.Equal(t, transfer.FromAccountID, args.FromAccountID)

		// Check amounts
		diff1 := account2.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account1.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		require.True(t, diff2 > 0)
		require.True(t, diff2%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
	}

}
