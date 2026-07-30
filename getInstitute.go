package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
)

func getInstitutePapers(mtid string) ([]Paper, error) {
	params := url.Values{}
	params.Add("cond", "institutes;inia;"+mtid)
	params.Add("cond", "published;eq;true")
	params.Add("cond", "core;eq;true")
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
		return nil, fmt.Errorf("institute %s: %w", mtid, err)
	}
	return getPapers(mtmtResponse, "-1"), nil
}

func getUnique(papers []Paper) []Paper {
	idmap := make(map[int]Paper)
	for _, paper := range papers {
		idmap[paper.Mtid] = paper
	}
	return lo.Values[int, Paper](idmap)
}

func getInstitutes(mtids []string) (PaperResponse, error) {
	type result struct {
		papers []Paper
		err    error
	}
	results := make(chan result, len(mtids))
	var wg sync.WaitGroup
	for _, id := range mtids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			papers, err := getInstitutePapers(id)
			results <- result{papers, err}
		}(id)
	}
	wg.Wait()
	close(results)

	var papers []Paper
	for r := range results {
		if r.err != nil {
			// One institute failing must not overwrite a good cache with a
			// partial list; the caller falls back to the existing cache.
			return PaperResponse{}, r.err
		}
		papers = append(papers, r.papers...)
	}
	papers = getUnique(papers)
	return PaperResponse{Papers: papers, Time: time.Now().Unix()}, nil
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
		return
	}
	if !isgoodparam {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - not an mtid"))
		return
	}
	sort.Strings(mtid)
	mtidstring := strings.Join(mtid, "_")
	log.Printf("/insitute %s\n", mtidstring)
	filename := "institutes_" + mtidstring + ".json"
	serveCached(w, filename, func() (PaperResponse, error) { return getInstitutes(mtid) })
}
