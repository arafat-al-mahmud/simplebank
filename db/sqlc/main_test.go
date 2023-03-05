package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var db *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}
	testQueries = New(db)
	os.Exit(m.Run())
}
