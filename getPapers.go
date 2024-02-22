package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getJournal(apistring string) string {
	req, err := url.Parse("https://m2.mtmt.hu/" + apistring)
	if err != nil {
		log.Fatalln(err)
		return ""
	}

	resp, err := http.Get(req.String())
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	journalResponse := JournalResponse{}
	err = json.Unmarshal([]byte(body), &journalResponse)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Title(strings.ToLower(journalResponse.Content.Title))
}

func getJournals(papers []Paper) []Paper {
	marshalledjson, _ := os.ReadFile("journalmap.json")
	journalmap := make(map[string]string)
	json.Unmarshal(marshalledjson, &journalmap)
	for i, paper := range papers {
		if title, ok := journalmap[paper.Journal]; ok {
			papers[i].Journal = title
		} else {
			title := getJournal(paper.Journal)
			journalmap[paper.Journal] = title
			papers[i].Journal = title
		}
	}
	marshalledjson, _ = json.Marshal(journalmap)
	_ = os.WriteFile("journalmap.json", marshalledjson, 0644)
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
			Sjr:                 content.Sjr,
		}
		papers = append(papers, paper)
	}

	return papers
}
