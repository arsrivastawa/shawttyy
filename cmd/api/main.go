package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/arsrivastawa/shawttyy/internal/api"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/service"
	"github.com/arsrivastawa/shawttyy/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
	baseURL       = "http://localhost:8080"
	defaultNodeID = 1
)

func runDBMigrations(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	err = migration.Up()

	if err != nil && err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	} else {
		log.Println("Database migrated successfully")
	}

}

func main() {

	_, filename, _, _ := runtime.Caller(0)

	baseDir := strings.Split(filepath.Dir(filename), "/")

	newPath := strings.Join(baseDir[:len(baseDir)-2], "/")

	fmt.Println(newPath)

	targetPath := filepath.Join(newPath, "internal/migrations", "000001_create_urls_table.up.sql")

	fmt.Println(targetPath)

	fmt.Println("The Shawttyy is up and running!!!")
	_ = godotenv.Load(".env")
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable is required but not set")
		os.Exit(1)
	}

	runDBMigrations("file://internal/migrations", connStr)

	db := storage.ConnectDB(connStr)
	defer db.Close()

	seq, err := sequencer.New(defaultNodeID)
	if err != nil {
		panic(err)
	}

	store := storage.NewPostgresStore(db)
	svc := service.New(store, seq)
	handler := api.NewHandler(svc, baseURL)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler)

	fmt.Printf("Listening on %s\n", baseURL)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
