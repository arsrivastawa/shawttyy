package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arsrivastawa/shawttyy/internal/api"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/service"
	"github.com/arsrivastawa/shawttyy/internal/storage"
	"github.com/joho/godotenv"
)

const (
	baseURL       = "http://localhost:8080"
	defaultNodeID = 1
)

func main() {
	fmt.Println("The Shawttyy is up and running!!!")
	_ = godotenv.Load(".env")
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable is required but not set")
		os.Exit(1)
	}

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
