package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbSource = "postgresql://postgres:jaw123@localhost:8892/curso?sslmode=disable"
)

var testQueries *Queries
var testConnPool *pgxpool.Pool
var testStore *Store

// Entry Point of all db tests
func TestMain(m *testing.M) {
	// Getting new connection pool
	testConnPool, err := pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("Could not connect to data source: ", err)
	}

	// Get the instance of queries for a specific database
	testQueries = New(testConnPool)

	// Get the isntance for Store
	testStore = NewStore(testConnPool)

	// Run tests
	os.Exit(m.Run())
}
