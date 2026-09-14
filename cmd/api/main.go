package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arsrivastawa/shawttyy/internal/api"
	"github.com/arsrivastawa/shawttyy/internal/cache"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/service"
	"github.com/arsrivastawa/shawttyy/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
	defaultNodeID      = 1
	defaultCacheExpiry = 1 * time.Hour
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
	fmt.Println("The Shawttyy is up and running!!!")
	_ = godotenv.Load(".env")
	connStr := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("SERVER_URL")
	redisURL := os.Getenv("REDIS_URL")

	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable is required but not set")
		os.Exit(1)
	}

	runDBMigrations("file://internal/migrations", connStr)

	db := storage.ConnectDB(connStr)
	rdb, ctx := cache.ConnectRedis(redisURL)
	defer db.Close()

	seq, err := sequencer.New(defaultNodeID)
	if err != nil {
		panic(err)
	}

	store := storage.NewPostgresStore(db)
	cache := cache.NewRedisCache(rdb, ctx, defaultCacheExpiry)
	svc := service.New(store, cache, seq)
	handler := api.NewHandler(svc, baseURL)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler)

	fmt.Printf("Listening on %s\n", baseURL)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
