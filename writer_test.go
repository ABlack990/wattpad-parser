package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteMetadataJsonFile(t *testing.T) {
	outputDir := t.TempDir()

	metadata := storyMetadata{
		Name:             "Duplicity [h.s]",
		URL:              "https://www.wattpad.com/story/138202651-duplicity-h-s",
		CommentCount:     32244839,
		Author:           storyPerson{Name: "julez!<3", URL: "https://www.wattpad.com/user/happydays1d"},
		IsFamilyFriendly: false,
	}

	if err := writeMetadataJsonFile(metadata, outputDir); err != nil {
		t.Fatalf("writeMetadataJsonFile returned an unexpected error: %v", err)
	}

	metadataFilePath := filepath.Join(outputDir, "metadata.json")
	fileBytes, err := os.ReadFile(metadataFilePath)
	if err != nil {
		t.Fatalf("failed to read metadata.json file: %v", err)
	}

	var written storyMetadata
	if err := json.Unmarshal(fileBytes, &written); err != nil {
		t.Fatalf("failed to unmarshal written metadata.json: %v", err)
	}

	if written != metadata {
		t.Errorf("written metadata does not match input.\ngot:  %+v\nwant: %+v", written, metadata)
	}
}

func TestWriteMetadataJsonFile_InvalidOutputDir(t *testing.T) {
	nonExistentDir := filepath.Join(t.TempDir(), "does-not-exist")

	if err := writeMetadataJsonFile(storyMetadata{}, nonExistentDir); err == nil {
		t.Fatal("expected an error when writing to a non-existent directory, got nil")
	}
}

func TestWriteChapter_Markdown(t *testing.T) {
	outputDir := t.TempDir()
	record := chapterRecord{title: "chapter-1", url: "https://www.wattpad.com/chapter-1"}

	if err := writeChapter("<p><i>Hello world</i></p>", record, outputDir, FormatMarkdown); err != nil {
		t.Fatalf("writeChapter returned an unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "chapter-1.md"))
	if err != nil {
		t.Fatalf("failed to read chapter-1.md: %v", err)
	}

	if !strings.Contains(string(content), "Hello world") {
		t.Errorf("expected markdown file to contain %q, got: %s", "Hello world", content)
	}
}

func TestWriteChapter_Html(t *testing.T) {
	outputDir := t.TempDir()
	record := chapterRecord{title: "chapter-1", url: "https://www.wattpad.com/chapter-1"}
	rawHtml := "<p><i>Hello world</i></p>"

	if err := writeChapter(rawHtml, record, outputDir, FormatHtml); err != nil {
		t.Fatalf("writeChapter returned an unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "chapter-1.html"))
	if err != nil {
		t.Fatalf("failed to read chapter-1.html: %v", err)
	}

	if string(content) != rawHtml {
		t.Errorf("expected chapter-1.html to contain the raw HTML unchanged.\ngot:  %s\nwant: %s", content, rawHtml)
	}
}

func TestWriteChapter_Text(t *testing.T) {
	outputDir := t.TempDir()
	record := chapterRecord{title: "chapter-1", url: "https://www.wattpad.com/chapter-1"}

	if err := writeChapter("<p><i>Hello world</i></p>", record, outputDir, FormatText); err != nil {
		t.Fatalf("writeChapter returned an unexpected error: %v", err)
	}

	// FormatText is not yet implemented, so no chapter file should be written
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("failed to read output directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected no files to be written for FormatText, got: %v", entries)
	}
}

func TestSanitizeFileName(t *testing.T) {
	if got := sanitizeFileName("normal-title"); got != "normal-title" {
		t.Errorf("expected sanitizeFileName to leave a normal title unchanged, got %q", got)
	}

	for _, invalid := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"} {
		input := "a" + invalid + "b"
		want := "a-b"
		if got := sanitizeFileName(input); got != want {
			t.Errorf("sanitizeFileName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestWriteChapter_TitleWithSlash(t *testing.T) {
	outputDir := t.TempDir()
	record := chapterRecord{title: "-read me / trailer ", url: "https://www.wattpad.com/chapter-1"}

	if err := writeChapter("<p><i>Hello world</i></p>", record, outputDir, FormatMarkdown); err != nil {
		t.Fatalf("writeChapter returned an unexpected error: %v", err)
	}

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("failed to read output directory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one file to be written, got: %v", entries)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "-read me - trailer .md"))
	if err != nil {
		t.Fatalf("failed to read sanitized chapter file: %v", err)
	}

	if !strings.Contains(string(content), "Hello world") {
		t.Errorf("expected markdown file to contain %q, got: %s", "Hello world", content)
	}
}

func TestWriteChapter_UnsupportedFormat(t *testing.T) {
	outputDir := t.TempDir()
	record := chapterRecord{title: "chapter-1", url: "https://www.wattpad.com/chapter-1"}

	err := writeChapter("<p>content</p>", record, outputDir, OutputFormat("bogus"))
	if err == nil {
		t.Fatal("expected an error for an unsupported output format, got nil")
	}
}
