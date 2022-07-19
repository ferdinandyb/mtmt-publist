package main

type Author struct {
	index      int
	familyName string
	givenName  string
	isUser     bool
}

type Paper struct {
	index               int
	journal             string
	year                int
	title               string
	doi                 string
	citation            int
	independentCitation int
	authors             []Author
}

type mtmtResponse struct {
	Content []struct {
		Title               string `json:"title"`
		Year                string `json:"publishedYear"`
		Citation            int    `json:"citationCount"`
		IndependentCitation int    `json:"independentCitationCount"`
		Authorships         []struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
			Type       struct {
				Mtid int `json:"mtid"`
			} `json:"type"`
		} `json:"authorships"`
		Identifiers []struct {
			RealUrl string `json:"realUrl"`
			Source  struct {
				Label string `json:"label"`
			} `json:"source"`
		} `json:"identifiers"`
	} `json:"content"`
}
