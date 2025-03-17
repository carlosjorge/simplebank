package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testPool, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	testQueries = New(testPool)

	os.Exit(m.Run())
}
