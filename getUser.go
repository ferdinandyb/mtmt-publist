package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getJournal(apistring string) string {
	return "hello"
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
	journals := make(map[string]string)
	for index, content := range mtmtResponse.Content {
		journals[content.Journal.Link] = ""
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
			fmt.Println(mtid_asint)
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
		}
		response = append(response, paper)
	}
	fmt.Println(journals)
	json.NewEncoder(w).Encode(response)

}

// var wg sync.WaitGroup
// journalsChan := make(chan string, len(content.Authorships))
// for i, author := range content.Authorships {
// 	wg.Add(1)
// 	i := i
// 	go func(author AuthorShip) {
// 		defer wg.Done()
// 		journalsChan <- getJournal()
// 	}(journal)
// }
