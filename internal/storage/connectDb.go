package storage

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(connStr string) *sql.DB {
	if connStr == "" {
		log.Fatal("Connection string is required but not set")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error connecting to the Database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error while pinging the database: ", err)
	}

	println("Successfully connected to the DB")
	return db
}
