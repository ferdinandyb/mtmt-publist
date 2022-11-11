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

func getJournals(papers []Paper) []Paper {
	var journals []string
	for _, paper := range papers {
		journals = append(journals, paper.Journal)
	}
	journals = lo.Uniq[string](journals)
	fmt.Println("starting on journals")
	journal_titles := lop.Map[string, string](journals, func(x string, _ int) string { return getJournal(x) })
	journalmap := make(map[string]string)
	for i := 0; i < len(journals); i++ {
		journalmap[journals[i]] = journal_titles[i]
	}
	fmt.Println("got papers")
	papers = lo.Map[Paper, Paper](papers, func(x Paper, _ int) Paper {
		x.Journal = journalmap[x.Journal]
		return x
	})
	return papers
}

func getPapers(mtmtResponse MtmtResponse, userMtid string) []Paper {
	var papers []Paper
	for index, content := range mtmtResponse.Content {
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
		mtid_asint, _ := strconv.Atoi(userMtid)
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
			Mtid:                content.Mtid,
			Index:               index,
			Title:               content.Title,
			Year:                content.Year,
			Citation:            content.Citation,
			IndependentCitation: content.IndependentCitation,
			Doi:                 doi,
			Authors:             authors,
			Journal:             content.Journal.Link,
		}
		papers = append(papers, paper)
	}

	return papers
}
