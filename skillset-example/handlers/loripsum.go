package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func Loripsum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Loripsum Called")

	params := &struct {
		NumberOfParagraphs int    `json:"number_of_paragraphs"`
		ParagraphLength    string `json:"paragraph_length"`
		Decorate           bool   `json:"decorate"`
		Link               bool   `json:"link"`
		UnorderedLists     bool   `json:"unordered_lists"`
		NumberedLists      bool   `json:"numbered_lists"`
		DescriptionLists   bool   `json:"description_lists"`
		BlockQuotes        bool   `json:"blockquotes"`
		Code               bool   `json:"code"`
		Headers            bool   `json:"headers"`
		AllCaps            bool   `json:"all_caps"`
		Prude              bool   `json:"prude"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := "api"
	if params.NumberOfParagraphs != 0 {
		if params.NumberOfParagraphs < 0 {
			p = path.Join(p, "1")
		} else if params.NumberOfParagraphs > 10 {
			p = path.Join(p, "10")
		} else {
			p = path.Join(p, fmt.Sprintf("%v", params.NumberOfParagraphs))
		}
	}

	if params.ParagraphLength == "short" || params.ParagraphLength == "medium" || params.ParagraphLength == "long" || params.ParagraphLength == "verylong" {
		p = path.Join(p, params.ParagraphLength)
	}

	if params.Decorate {
		p = path.Join(p, "decorate")
	}

	if params.Link {
		p = path.Join(p, "link")
	}

	if params.UnorderedLists {
		p = path.Join(p, "ul")
	}

	if params.NumberedLists {
		p = path.Join(p, "ol")
	}

	if params.DescriptionLists {
		p = path.Join(p, "dl")
	}

	if params.BlockQuotes {
		p = path.Join(p, "bq")
	}

	if params.Code {
		p = path.Join(p, "code")
	}

	if params.Headers {
		p = path.Join(p, "headers")
	}

	if params.AllCaps {
		p = path.Join(p, "allcaps")
	}

	if params.Prude {
		p = path.Join(p, "prude")
	}

	u, _ := url.Parse("https://loripsum.net")
	u.Path = p

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Println("Something unrecoverable went wrong")
		return
	}
}
