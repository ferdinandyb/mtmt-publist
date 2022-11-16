package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
)

var CACHETIME int64

func main() {
	CACHETIME = 60 * 60 * 24
	file, logerr := os.OpenFile("mtmt-publist.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if logerr != nil {
		log.Fatal(logerr)
	}

	log.SetOutput(file)

	mux := http.NewServeMux()
	mux.HandleFunc("/user", handleGetUser)
	mux.HandleFunc("/institute", handleGetInstitute)
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("alive")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - alive"))
	},
	)
	var port string
	flag.StringVar(&port, "port", "3333", "specify port")
	flag.Parse()
	servererr := http.ListenAndServe(":"+port, mux)
	if errors.Is(servererr, http.ErrServerClosed) {
		log.Println("server closed")
	} else if servererr != nil {
		log.Printf("error starting server: %s\n", servererr)
		os.Exit(1)
	}
}
