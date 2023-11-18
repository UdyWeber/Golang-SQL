package db

import (
	"awesomeProject/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    utils.GenerateRandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	accountFromDb, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFromDb)

	require.Equal(t, account.ID, accountFromDb.ID)
	require.Equal(t, account.Owner, accountFromDb.Owner)
	require.Equal(t, account.Currency, accountFromDb.Currency)
	require.Equal(t, account.Balance, accountFromDb.Balance)
	require.WithinDuration(t, account.CreatedAt, accountFromDb.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	args := UpdateAccountParams{
		NewBalance:  utils.RandomMoney(),
		NewOwner:    utils.GenerateRandomOwner(),
		NewCurrency: utils.RandomCurrency(),
		ID:          account.ID,
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, accountUpdated.Owner, args.NewOwner)
	require.Equal(t, accountUpdated.Currency, args.NewCurrency)
	require.NotEqual(t, account.Balance, accountUpdated.Balance)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)

	require.Equal(t, accountUpdated.Balance, args.NewBalance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	account, err = testQueries.GetAccount(context.Background(), account.ID)

	require.Empty(t, account)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
