package api

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hey, The Shawttyy is up and running!!!")

	w.WriteHeader(http.StatusOK)

}
