package db

import (
	"awesomeProject/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, args.Amount)
	require.Equal(t, entry.AccountID, account.ID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	entryFromDB, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)

	require.Equal(t, entry.ID, entryFromDB.ID)
	require.Equal(t, entry.AccountID, entryFromDB.AccountID)
	require.Equal(t, entry.Amount, entryFromDB.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryFromDB.CreatedAt, time.Second)
}

func TestGetEntriesByAccount(t *testing.T) {
	var lastEntry Entry

	for i := 0; i < 2; i++ {
		lastEntry = createRandomEntry(t)
	}

	entries, err := testQueries.GetEntriesByAccount(context.Background(), lastEntry.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 1)
	require.Contains(t, entries, lastEntry)

	entry := entries[0]

	require.Equal(t, entry.AccountID, lastEntry.AccountID)
	require.Equal(t, entry.ID, lastEntry.ID)
	require.Equal(t, entry.Amount, lastEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt, lastEntry.CreatedAt, time.Second)
}

func TestGetAllEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	args := GetAllEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetAllEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entries := range entries {
		require.NotEmpty(t, entries)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)
	anotherAccount := createRandomAccount(t)

	args := UpdateEntryParams{
		Amount:    utils.RandomMoney(),
		AccountID: anotherAccount.ID,
		ID:        entry.ID,
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)

	require.Equal(t, updatedEntry.Amount, args.Amount)
	require.Equal(t, updatedEntry.AccountID, anotherAccount.ID)
	require.Equal(t, updatedEntry.ID, entry.ID)
	require.WithinDuration(t, entry.CreatedAt, updatedEntry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	entry, err = testQueries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, entry)

}
