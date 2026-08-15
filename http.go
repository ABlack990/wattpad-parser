package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpConfig struct {
	Timeout time.Duration
}

var httpConfig = HttpConfig{Timeout: 10 * time.Second}

func fetchRawHtml(fetchUrl string) (string, error) {
	client := &http.Client{Timeout: httpConfig.Timeout}

	req, err := http.NewRequest(http.MethodGet, fetchUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
