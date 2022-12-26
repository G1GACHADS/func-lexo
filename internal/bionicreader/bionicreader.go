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
)

func Convert(content string, fixation, saccade int) (string, error) {
	payload := bytes.NewBufferString(url.Values{
		"content":       {content},
		"fixation":      {strconv.FormatInt(int64(fixation), 10)},
		"saccade":       {strconv.FormatInt(int64(saccade), 10)},
		"response_type": {"html"},
		"request_type":  {"html"},
	}.Encode())

	req, err := http.NewRequest(http.MethodPost, os.Getenv("BIONIC_READING_API_URL"), payload)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Set("X-RapidAPI-Host", os.Getenv("RAPID_API_HOST"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	convertedText, err := parseConvertedText(resp.Body)
	if err != nil {
		return "", err
	}

	return convertedText, nil
}

func parseConvertedText(body io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	html, err := doc.Find(".bionic-reader .bionic-reader-content .bionic-reader-container").First().Html()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(html), nil
}
