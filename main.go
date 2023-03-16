package main

import (
	"database/sql"
	"log"

	"github.com/arafat-al-mahmud/simplebank/api"
	db "github.com/arafat-al-mahmud/simplebank/db/sqlc"
	"github.com/arafat-al-mahmud/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't read config file, ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Error starting server : ", err)
	}
}
