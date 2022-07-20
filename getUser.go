package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getJournal(apistring string) string {
	req, err := url.Parse("https://m2.mtmt.hu/" + apistring)
	if err != nil {
		return ""
	}

	resp, err := http.Get(req.String())
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	journalResponse := JournalResponse{}
	err = json.Unmarshal([]byte(body), &journalResponse)
	return strings.Title(strings.ToLower(journalResponse.Content.Title))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /user request\n")
	mtid := r.URL.Query().Get("mtid")
	base, err := url.Parse("https://m2.mtmt.hu/api/publication")
	if err != nil {
		return
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
	var response []Paper
	var journals []string
	for index, content := range mtmtResponse.Content {
		journals = append(journals, content.Journal.Link)
		var doi string
		for _, identifier := range content.Identifiers {
			if identifier.Source.Label == "DOI" {
				doi = identifier.RealUrl
				if doi == "" {
					doi = "https://doi.org/" + strings.Split(identifier.Label, " ")[0]
				}
			}
		}
		var authors []Author
		mtid_asint, _ := strconv.Atoi(mtid)
		for author_i, author := range content.Authorships {
			if author.Type.Mtid == 1 {
				authors = append(authors, Author{
					Index:      author_i,
					FamilyName: author.FamilyName,
					GivenName:  author.GivenName,
					IsUser:     author.Author.Mtid == mtid_asint,
				})
			}
		}
		paper := Paper{
			Index:               index,
			Title:               content.Title,
			Year:                content.Year,
			Citation:            content.Citation,
			IndependentCitation: content.IndependentCitation,
			Doi:                 doi,
			Authors:             authors,
			Journal:             content.Journal.Link,
		}
		response = append(response, paper)
	}

	journals = lo.Uniq[string](journals)
	journal_titles := lop.Map[string, string](journals, func(x string, _ int) string { return getJournal(x) })
	journalmap := make(map[string]string)
	for i := 0; i < len(journals); i++ {
		journalmap[journals[i]] = journal_titles[i]
	}
	response = lo.Map[Paper, Paper](response, func(x Paper, _ int) Paper {
		x.Journal = journalmap[x.Journal]
		return x
	})
	json.NewEncoder(w).Encode(response)

}
