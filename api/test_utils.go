package api

import (
	"io"
	"net/http"
)

func NewRequestWithTestHeader(method, url string, body io.Reader, header map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k := range header {
		req.Header.Add(k, header[k])
	}

	return req, nil
}
