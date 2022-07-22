package main

import (
	"encoding/json"
	"github.com/samber/lo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func getInstitutePapers(mtid string) ([]Paper, error) {

	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return make([]Paper, 0), err
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
		return make([]Paper, 0), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]Paper, 0), err
	}
	mtmtResponse := MtmtResponse{}
	err = json.Unmarshal([]byte(body), &mtmtResponse)
	papers := getPapers(mtmtResponse, "-1")
	return papers, nil
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
	for _, id := range mtids {
		inst_papers, err := getInstitutePapers(id)
		if err != nil {
			return PaperResponse{}, err
		}
		papers = append(papers, inst_papers...)
	}
	papers = getUnique(papers)

	retval := PaperResponse{Papers: papers, Time: time.Now().Unix()}
	return retval, nil

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
					jsonresp, _ = ioutil.ReadFile(filename)
					w.Write(jsonresp)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - MTMT is probably not available and no fallback exists"))
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

}
