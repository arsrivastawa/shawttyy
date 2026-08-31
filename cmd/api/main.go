package main

import (
	"fmt"
	"net/http"

	"github.com/arsrivastawa/shawttyy/internal/api"
)

func main() {
	fmt.Println("The Shawttyy is up and running!!!")

	mux := http.NewServeMux()

	api.RegisterHelloRouter(mux)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

}
