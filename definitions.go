package main

type Author struct {
	index      int
	familyName string
	givenName  string
	isUser     bool
}

type Paper struct {
	Index               int
	Journal             string
	Year                int
	Title               string
	Doi                 string
	Citation            int
	IndependentCitation int
	Authors             []Author
}

type MtmtResponse struct {
	Content []struct {
		Title               string `json:"title"`
		Year                int    `json:"publishedYear"`
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
