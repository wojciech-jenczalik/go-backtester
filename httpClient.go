package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func get(url string) ([]byte, error) {
	resp, err := http.Get(prepURL(url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func prepURL(url string) string {

	var queryChar byte
	if strings.Contains(url, "?") {
		queryChar = '&'
	} else {
		queryChar = '?'
	}

	var sb strings.Builder
	sb.WriteString(FmpUrl)
	sb.WriteString(url)
	sb.WriteByte(queryChar)
	sb.WriteString(FmpApiKeyKey)
	sb.WriteByte('=')
	sb.WriteString(os.Getenv("FMP_API_KEY"))

	return sb.String()
}
