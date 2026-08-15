package main

import (
	"os"
	"strings"
	"testing"
)

func TestConvertHtmlToMarkdown_Duplicity00Content(t *testing.T) {
	rawHtml, err := os.ReadFile("test_files/duplicity-00-content.html")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	markdown, err := convertHtmlToMarkdown(string(rawHtml))
	if err != nil {
		t.Fatalf("convertHtmlToMarkdown returned an unexpected error: %v", err)
	}

	if markdown == "" {
		t.Fatal("expected non-empty markdown")
	}

	if strings.Contains(markdown, "<") || strings.Contains(markdown, ">") {
		t.Errorf("expected markdown to contain no raw HTML tags, got: %s", markdown)
	}

	if !strings.Contains(markdown, "I gasped in horror") {
		t.Errorf("expected markdown to contain known chapter text, got: %s", markdown)
	}

	if strings.Contains(markdown, "trinityAudioPlaceholder") {
		t.Errorf("expected trinityAudioPlaceholder div to not appear in markdown, got: %s", markdown)
	}

	// smart-quote decoded HTML entities should survive the conversion
	if !strings.Contains(markdown, "It seems we have a real problem now") {
		t.Errorf("expected markdown to contain decoded quoted text, got: %s", markdown)
	}
}

func TestConvertHtmlToMarkdown_EmptyInput(t *testing.T) {
	markdown, err := convertHtmlToMarkdown("")
	if err != nil {
		t.Fatalf("convertHtmlToMarkdown returned an unexpected error: %v", err)
	}

	if markdown != "" {
		t.Errorf("expected empty markdown for empty input, got: %q", markdown)
	}
}

func TestConvertHtmlToMarkdown_BasicFormatting(t *testing.T) {
	html := `<p><b><i>Bold and italic</i></b></p><p><i>Just italic</i></p>`

	markdown, err := convertHtmlToMarkdown(html)
	if err != nil {
		t.Fatalf("convertHtmlToMarkdown returned an unexpected error: %v", err)
	}

	if !strings.Contains(markdown, "Bold and italic") {
		t.Errorf("expected markdown to contain %q, got: %s", "Bold and italic", markdown)
	}

	if !strings.Contains(markdown, "Just italic") {
		t.Errorf("expected markdown to contain %q, got: %s", "Just italic", markdown)
	}

	if !strings.Contains(markdown, "*") {
		t.Errorf("expected markdown to retain emphasis markers, got: %s", markdown)
	}
}
