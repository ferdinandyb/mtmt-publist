package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getUser(mtid string) PaperResponse {
	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return PaperResponse{}
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
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	mtmtResponse := MtmtResponse{}
	err = json.Unmarshal([]byte(body), &mtmtResponse)
	papers := getPapers(mtmtResponse, mtid)
	retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
	return retval

}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /user request\n")
	mtid := r.URL.Query().Get("mtid")
	filename := "user_" + mtid + ".json"
	info, err := os.Stat(filename)
	var jsonresp []byte
	if err != nil || time.Now().Unix()-info.ModTime().Unix() > CACHETIME {
		response := getUser(mtid)
		jsonresp, _ = json.Marshal(response)
		_ = ioutil.WriteFile(filename, jsonresp, 0644)
	} else {
		jsonresp, _ = ioutil.ReadFile(filename)
	}
	w.Write(jsonresp)
}
