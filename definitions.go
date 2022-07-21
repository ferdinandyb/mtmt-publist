package main

type Author struct {
	Index      int
	FamilyName string
	GivenName  string
	IsUser     bool
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
	Mtid                int
}

type PaperResponse struct {
	Papers []Paper
	Time   int64
}

type JournalResponse struct {
	Content struct {
		Title string `json:"title"`
	} `json:"content"`
}

type AuthorShip struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	Author     struct {
		Mtid int `json:"mtid"`
	} `json:"author"`
	Type struct {
		Mtid int `json:"mtid"`
	} `json:"type"`
}

type MtmtResponse struct {
	Content []struct {
		Mtid                int          `json:"mtid"`
		Title               string       `json:"title"`
		Year                int          `json:"publishedYear"`
		Citation            int          `json:"citationCount"`
		IndependentCitation int          `json:"independentCitationCount"`
		Authorships         []AuthorShip `json:"authorships"`
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
	} `json:"content"`
}
