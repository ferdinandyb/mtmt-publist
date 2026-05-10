package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// fetchAllPages fetches all pages from the MTMT API for the given base params
// and returns a single MtmtResponse with all content entries merged.
// The caller should set all cond/sort/fields params; size and page are managed here.
func fetchAllPages(params url.Values) (MtmtResponse, error) {
	params.Set("size", "5000")
	var allContent []struct {
		Mtid                int          `json:"mtid"`
		Title               string       `json:"title"`
		Year                int          `json:"publishedYear"`
		Citation            int          `json:"citationCount"`
		IndependentCitation int          `json:"independentCitationCount"`
		Authorships         []AuthorShip `json:"authorships"`
		Sjr                 string       `json:"ratingsForSort"`
		Identifiers         []struct {
			RealUrl string `json:"realUrl"`
			Label   string `json:"label"`
			Source  struct {
				Label string `json:"label"`
			} `json:"source"`
		} `json:"identifiers"`
		Journal struct {
			Link string `json:"link"`
		} `json:"journal"`
	}

	for page := 1; ; page++ {
		params.Set("page", strconv.Itoa(page))
		base, err := url.Parse("https://m2.mtmt.hu/api/publication")
		if err != nil {
			return MtmtResponse{}, err
		}
		base.RawQuery = params.Encode()

		resp, err := http.Get(base.String())
		if err != nil {
			return MtmtResponse{}, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return MtmtResponse{}, err
		}

		var pageResp MtmtResponse
		if err := json.Unmarshal(body, &pageResp); err != nil {
			log.Printf("fetchAllPages: failed to unmarshal page %d: %v\nbody: %.200s", page, err, body)
			return MtmtResponse{}, fmt.Errorf("unmarshal page %d: %w", page, err)
		}

		allContent = append(allContent, pageResp.Content...)

		if pageResp.Paging.Last {
			break
		}
	}

	return MtmtResponse{Content: allContent}, nil
}

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
