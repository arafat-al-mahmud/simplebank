package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/arafat-al-mahmud/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var db *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Couldn't read config file, ", err)
	}
	db, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}
	testQueries = New(db)
	os.Exit(m.Run())
}
