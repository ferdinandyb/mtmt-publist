package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var CACHETIME int64

// serveCached writes filename to w if it is fresh, otherwise calls fetch,
// caches the result to filename, and writes that. If fetch fails, a still
// existing (stale) cache file is served instead of a 500.
func serveCached(w http.ResponseWriter, filename string, fetch func() (PaperResponse, error)) {
	info, fileerr := os.Stat(filename)
	if fileerr == nil && time.Now().Unix()-info.ModTime().Unix() < CACHETIME {
		jsonresp, _ := os.ReadFile(filename)
		w.Write(jsonresp)
		return
	}

	response, err := fetch()
	if err != nil {
		if fileerr == nil {
			jsonresp, _ := os.ReadFile(filename)
			w.Write(jsonresp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - MTMT is probably not available and no fallback exists"))
		return
	}

	jsonresp, _ := json.Marshal(response)
	w.Write(jsonresp)
	_ = os.WriteFile(filename, jsonresp, 0644)
}

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
