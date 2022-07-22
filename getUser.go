package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getUser(mtid string) (PaperResponse, error) {
	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return PaperResponse{}, err
	}

	// Query params
	params := url.Values{}
	params.Add("cond", "published;eq;true")
	params.Add("cond", "core;eq;true")
	params.Add("cond", "authors.mtid;eq;"+mtid)
	params.Add("cond", "category.mtid;eq;1")
	params.Add("cond", "type.mtid;eq;24")
	params.Add("cond", "languages.label;eq;Angol")
	params.Add("sort", "publishedYear,desc")
	params.Add("sort", "firstAuthor,asc")
	params.Add("size", "10000")
	params.Add("size", "10000")
	params.Add("fields", "template")
	params.Add("labelLang", "hun")
	params.Add("cite_type", "2")
	params.Add("page", "1")
	params.Add("format", "json")
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return PaperResponse{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PaperResponse{}, err
	}
	mtmtResponse := MtmtResponse{}
	err = json.Unmarshal([]byte(body), &mtmtResponse)
	papers := getPapers(mtmtResponse, mtid)
	retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
	return retval, nil

}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /user request\n")
	mtid := r.URL.Query().Get("mtid")
	filename := "user_" + mtid + ".json"
	info, fileerr := os.Stat(filename)
	var jsonresp []byte
	if fileerr != nil || time.Now().Unix()-info.ModTime().Unix() >= CACHETIME {
		response, err := getUser(mtid)
		if err != nil {
			if fileerr == nil {
				jsonresp, _ = ioutil.ReadFile(filename)
				w.Write(jsonresp)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened!"))
			}
		} else {
			jsonresp, _ = json.Marshal(response)
			w.Write(jsonresp)
			_ = ioutil.WriteFile(filename, jsonresp, 0644)
		}
	} else {
		jsonresp, _ = ioutil.ReadFile(filename)
		w.Write(jsonresp)
	}
}
