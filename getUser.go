package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func getUser(mtid string) (PaperResponse, error) {
	params := url.Values{}
	params.Add("cond", "published;eq;true")
	params.Add("cond", "core;eq;true")
	params.Add("cond", "authors.mtid;eq;"+mtid)
	params.Add("cond", "category.mtid;eq;1")
	params.Add("cond", "type.mtid;eq;24")
	params.Add("cond", "languages.label;eq;Angol")
	params.Add("sort", "publishedYear,desc")
	params.Add("sort", "firstAuthor,asc")
	params.Add("fields", "template")
	params.Add("labelLang", "hun")
	params.Add("cite_type", "2")
	params.Add("format", "json")

	mtmtResponse, err := fetchAllPages(params)
	if err != nil {
		return PaperResponse{}, err
	}
	papers := getPapers(mtmtResponse, mtid)
	retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
	return retval, nil
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	mtid := r.URL.Query().Get("mtid")
	log.Printf("/user %s\n", mtid)
	regres, _ := regexp.MatchString(`^\d+$`, mtid)
	if mtid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - no MTID given"))
		return
	}
	if !regres {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - not an mtid"))
		return
	}
	serveCached(w, "user_"+mtid+".json", func() (PaperResponse, error) { return getUser(mtid) })
}
