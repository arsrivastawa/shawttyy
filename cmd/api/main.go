package main

import (
	"fmt"
	"net/http"

	"github.com/arsrivastawa/shawttyy/internal/api"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/service"
	"github.com/arsrivastawa/shawttyy/internal/storage"
)

const (
	baseURL       = "http://localhost:8080"
	defaultNodeID = 1
)

func main() {
	fmt.Println("The Shawttyy is up and running!!!")

	seq, err := sequencer.New(defaultNodeID)
	if err != nil {
		panic(err)
	}

	store := storage.NewInMemoryStore()
	svc := service.New(store, seq)
	handler := api.NewHandler(svc, baseURL)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler)

	fmt.Printf("Listening on %s\n", baseURL)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
