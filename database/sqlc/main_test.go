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

func TestMain(m *testing.M) {
	conn, err := sql.Open(databaseDriver, databaseSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}
	testQuery = New(conn)
	os.Exit(m.Run())
}