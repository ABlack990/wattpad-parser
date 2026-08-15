package main

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestExtractChapterURLs_DuplicityHome(t *testing.T) {
	rawHtml, err := os.ReadFile("test_files/duplicity-home.html")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	chapterRecords, err := extractChapterURLs(string(rawHtml))
	if err != nil {
		t.Fatalf("extractChapterURLs returned an unexpected error: %v", err)
	}

	if len(chapterRecords) == 0 {
		t.Errorf("expected at least one chapter URL, got 0")
	}

	for _, chapterRecord := range chapterRecords {
		title := chapterRecord.title
		chapterURL := chapterRecord.url
		if title == "" {
			t.Errorf("chapter title is empty for URL: %s", chapterURL)
		}
		if chapterURL == "" {
			t.Errorf("chapter URL is empty for title: %s", title)
		}

		if _, err := url.Parse(chapterURL); err != nil {
			t.Errorf("chapter URL %q for title %q is not a valid URL: %v", chapterURL, title, err)
		}

		if !strings.HasPrefix(chapterURL, "https://www.wattpad.com") {
			t.Errorf("chapter URL %q for title %q does not start with https://www.wattpad.com", chapterURL, title)
		}
	}

	if len(chapterRecords) != 102 {
		t.Errorf("expected 102 chapter URLs, got %d", len(chapterRecords))
	}
}

func TestExtractStoryMetadata_DuplicityHome(t *testing.T) {
	rawHtml, err := os.ReadFile("test_files/duplicity-home.html")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	metadata, err := extractStoryMetadata(string(rawHtml))
	if err != nil {
		t.Fatalf("extractStoryMetadata returned an unexpected error: %v", err)
	}

	if metadata.Name != "Duplicity [h.s]" {
		t.Errorf("expected name %q, got %q", "Duplicity [h.s]", metadata.Name)
	}
	if metadata.Author.Name == "" {
		t.Error("expected non-empty author name")
	}
	if metadata.URL != "https://www.wattpad.com/story/138202651-duplicity-h-s" {
		t.Errorf("unexpected url: %q", metadata.URL)
	}
	if metadata.CommentCount != 32244839 {
		t.Errorf("expected commentCount 32244839, got %d", metadata.CommentCount)
	}
	if metadata.InteractionStatistic != 117973435 {
		t.Errorf("expected interactionStatistic 117973435, got %d", metadata.InteractionStatistic)
	}
	if metadata.Publisher.Name != "Wattpad" {
		t.Errorf("expected publisher name %q, got %q", "Wattpad", metadata.Publisher.Name)
	}
	if metadata.DatePublished != "2018-02-11" {
		t.Errorf("expected datePublished %q, got %q", "2018-02-11", metadata.DatePublished)
	}
	if metadata.IsFamilyFriendly {
		t.Error("expected isFamilyFriendly to be false")
	}
}

func TestExtractChapterTextHtml_Duplicity00(t *testing.T) {
	rawHtml, err := os.ReadFile("test_files/duplicity-00.html")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	chapterHtml, err := extractChapterTextHtml(string(rawHtml))
	if err != nil {
		t.Fatalf("extractChapterTextHtml returned an unexpected error: %v", err)
	}

	if chapterHtml == "" {
		t.Fatal("expected non-empty chapter HTML")
	}

	if !strings.Contains(chapterHtml, "<p") {
		t.Errorf("expected chapter HTML to contain <p> tags, got: %s", chapterHtml)
	}

	if !strings.Contains(chapterHtml, "I was exuding out of total fear") {
		t.Errorf("expected chapter HTML to contain known chapter text, got: %s", chapterHtml)
	}

	if strings.Contains(chapterHtml, "trinityAudioPlaceholder") {
		t.Errorf("expected trinityAudioPlaceholder div to be removed, got: %s", chapterHtml)
	}
}
