package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var CACHETIME int64

func main() {
	CACHETIME = 60 * 60 * 24
	mux := http.NewServeMux()
	mux.HandleFunc("/user", handleGetUser)
	mux.HandleFunc("/institute", handleGetInstitute)

	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
