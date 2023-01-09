package bionicreader

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattn/godown"
)

type BionicReaderResult struct {
	Markdown string `json:"markdown"`
	Text     string `json:"text"`
}

func Convert(content string, fixation, saccade int) (BionicReaderResult, error) {
	payload := bytes.NewBufferString(url.Values{
		"content":       {content},
		"fixation":      {strconv.FormatInt(int64(fixation), 10)},
		"saccade":       {strconv.FormatInt(int64(saccade), 10)},
		"response_type": {"html"},
		"request_type":  {"html"},
	}.Encode())

	req, err := http.NewRequest(http.MethodPost, os.Getenv("BIONIC_READING_API_URL"), payload)
	if err != nil {
		return BionicReaderResult{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Set("X-RapidAPI-Host", os.Getenv("RAPID_API_HOST"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BionicReaderResult{}, err
	}
	defer resp.Body.Close()

	return parseConvertedText(resp.Body)
}

func parseConvertedText(body io.ReadCloser) (BionicReaderResult, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return BionicReaderResult{}, err
	}

	// Remove class attributes for each bold tags
	doc.Find(".bionic-reader-container").Find("b").RemoveAttr("class")

	doc.Find(".bionic-reader-container").Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#comment" {
			s.Remove()
		}
	})

	html, err := doc.Find(".bionic-reader-container").Html()
	if err != nil {
		return BionicReaderResult{}, err
	}

	html = strings.TrimSpace(html)

	var md bytes.Buffer
	if err := godown.Convert(&md, strings.NewReader(html), nil); err != nil {
		return BionicReaderResult{}, err
	}

	raw := doc.Find(".bionic-reader-container").Text()

	return BionicReaderResult{
		Markdown: md.String(),
		Text:     raw,
	}, nil
}
