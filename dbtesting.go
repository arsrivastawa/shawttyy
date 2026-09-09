package main

import (
	"database/sql"
	"log"

	// The underscore registers the driver with database/sql behind the scenes
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	connStr := ""

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error connecting to the Database: ", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error while pinging the database: ", err)
	}

	println("Successfully connected to the DB")
}
