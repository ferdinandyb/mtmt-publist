package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
)

func getInsitutePapers(mtid string) []Paper {

	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return make([]Paper, 0)
	}

	// Query params
	params := url.Values{}
	params.Add("cond", "institutes;inia;"+mtid)
	params.Add("cond", "published;eq;true")
	params.Add("cond", "core;eq;true")
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
	papers := getPapers(mtmtResponse, "-1")
	return papers
}

func getUnique(papers []Paper) []Paper {
	idmap := make(map[int]Paper)
	for _, paper := range papers {
		idmap[paper.Mtid] = paper
	}
	return lo.Values[int, Paper](idmap)

}

func getInstitutes(mtids []string) PaperResponse {
	var papers []Paper
	for _, id := range mtids {
		papers = append(papers, getInsitutePapers(id)...)
	}
	papers = getUnique(papers)

	retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
	return retval

}

func handleGetInstitute(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /insitute request\n")
	mtid := r.URL.Query()["mtid"]
	sort.Strings(mtid)
	filename := "institutes_" + strings.Join(mtid, "_") + ".json"
	info, err := os.Stat(filename)
	var jsonresp []byte
	if err != nil || time.Now().Unix()-info.ModTime().Unix() > CACHETIME {
		response := getInstitutes(mtid)
		jsonresp, _ = json.Marshal(response)
		_ = ioutil.WriteFile(filename, jsonresp, 0644)
	} else {
		jsonresp, _ = ioutil.ReadFile(filename)
	}
	w.Write(jsonresp)

}
