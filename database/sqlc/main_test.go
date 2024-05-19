package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	databaseDriver = "postgres"
	databaseSource = "postgresql://root:secretpassword@localhost:5432/simplebank?sslmode=disable"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(databaseDriver, databaseSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
