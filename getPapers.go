package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// issnSuffix matches the ISSN(s) MTMT appends to journal.label, e.g.
// "JOURNAL OF EXPERIMENTAL BIOLOGY 0022-0949 1477-9145".
var issnSuffix = regexp.MustCompile(`(\s+\d{4}-\d{3}[\dXx])+\s*$`)

// journalTitle derives a clean journal name from journal.label, avoiding an
// extra /api/journal request per publication.
func journalTitle(label string) string {
	return strings.Title(strings.ToLower(issnSuffix.ReplaceAllString(label, "")))
}

// fetchAllPages fetches all pages from the MTMT API for the given base params
// and returns a single MtmtResponse with all content entries merged.
// The caller should set all cond/sort/fields params; size and page are managed here.
func fetchAllPages(params url.Values) (MtmtResponse, error) {
	params.Set("size", "5000")
	var allContent []Publication

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
			Journal:             journalTitle(content.Journal.Label),
			Sjr:                 content.Sjr,
		}
		papers = append(papers, paper)
	}

	return papers
}
