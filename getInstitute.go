package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
)

func getInstitutePapers(mtid string, paperchan chan []Paper) {
	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return
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
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	mtmtResponse := MtmtResponse{}
	err = json.Unmarshal([]byte(body), &mtmtResponse)
	papers := getPapers(mtmtResponse, "-1")
	paperchan <- papers
	return
}

func getUnique(papers []Paper) []Paper {
	idmap := make(map[int]Paper)
	for _, paper := range papers {
		idmap[paper.Mtid] = paper
	}
	return lo.Values[int, Paper](idmap)
}

func getInstitutes(mtids []string) (PaperResponse, error) {
	var papers []Paper
	var wg sync.WaitGroup
	paperchan := make(chan []Paper)
	responsechan := make(chan PaperResponse)
	for _, id := range mtids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			getInstitutePapers(id, paperchan)
		}(id)
	}
	go func(responsechan chan PaperResponse) {
		for inst_papers := range paperchan {
			papers = append(papers, inst_papers...)
		}
		papers = getUnique(papers)
		papers = getJournals(papers)
		retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
		responsechan <- retval
	}(responsechan)
	wg.Wait()
	close(paperchan)

	return <-responsechan, nil
}

func handleGetInstitute(w http.ResponseWriter, r *http.Request) {
	mtid := r.URL.Query()["mtid"]
	isgoodparam := true
	for _, id := range mtid {
		regres, _ := regexp.MatchString(`^\d+$`, id)
		if !regres {
			isgoodparam = false
		}
	}
	if len(mtid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - no MTID given"))
	} else if !isgoodparam {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - not an mtid"))
	} else {
		sort.Strings(mtid)
		mtidstring := strings.Join(mtid, "_")
		log.Printf("/insitute %s\n", mtidstring)
		filename := "institutes_" + mtidstring + ".json"
		info, fileerr := os.Stat(filename)
		var jsonresp []byte
		if fileerr != nil || time.Now().Unix()-info.ModTime().Unix() >= CACHETIME {
			response, err := getInstitutes(mtid)
			if err != nil {
				if fileerr == nil {
					jsonresp, _ = os.ReadFile(filename)
					w.Write(jsonresp)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - MTMT is probably not available and no fallback exists"))
				}
			} else {
				jsonresp, _ = json.Marshal(response)
				w.Write(jsonresp)
				_ = os.WriteFile(filename, jsonresp, 0644)
			}
		} else {
			jsonresp, _ = os.ReadFile(filename)
			w.Write(jsonresp)
		}
	}
}
