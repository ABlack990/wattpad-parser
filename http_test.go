package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchRawHtml_Success(t *testing.T) {
	const wantBody = "<html><body>hello</body></html>"

	// httptest.NewServer spins up a real local HTTP server we can point fetchRawHtml at,
	// so no actual network call to wattpad.com is made.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(wantBody))
	}))
	defer server.Close()

	got, err := fetchRawHtml(server.URL)
	if err != nil {
		t.Fatalf("fetchRawHtml returned an unexpected error: %v", err)
	}

	if got != wantBody {
		t.Errorf("fetchRawHtml() = %q, want %q", got, wantBody)
	}
}

func TestFetchRawHtml_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchRawHtml(server.URL)
	if err == nil {
		t.Fatal("expected fetchRawHtml to return an error for a non-200 status, got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to mention status code 404, got: %v", err)
	}
}

func TestFetchRawHtml_InvalidUrl(t *testing.T) {
	_, err := fetchRawHtml("://not-a-valid-url")
	if err == nil {
		t.Fatal("expected fetchRawHtml to return an error for an invalid URL, got nil")
	}
}
