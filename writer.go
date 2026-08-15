package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Write metadata.json file
func writeMetadataJsonFile(metadata storyMetadata, outputDir string) error {
	metadataFilePath := fmt.Sprintf("%s/metadata.json", outputDir)

	file, err := os.Create(metadataFilePath)

	if err != nil {
		return fmt.Errorf("failed to create metadata.json file: %v", err)
	}
	defer file.Close()

	jsonBytesIndented, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to indented JSON: %v", err)
	}

	if _, err := file.Write(jsonBytesIndented); err != nil {
		return fmt.Errorf("failed to write metadata.json file: %v", err)
	}

	return nil
}

// Write chapter directory
func writeChapter(chapterRawHtml string, chapterRecord chapterRecord, outputDir string, outputFormat OutputFormat) error {
	title := sanitizeFileName(chapterRecord.title)

	switch outputFormat {
	case FormatMarkdown:
		content, err := convertHtmlToMarkdown(chapterRawHtml)
		if err != nil {
			return fmt.Errorf("failed to convert chapter HTML to Markdown: %v", err)
		}
		chapterFilePath := fmt.Sprintf("%s/%s.md", outputDir, title)
		return writeStringToFile(content, chapterFilePath)
	case FormatHtml:
		return writeStringToFile(chapterRawHtml, fmt.Sprintf("%s/%s.html", outputDir, title))
	case FormatText:
	default:
		return fmt.Errorf("unsupported output format D: : %s", outputFormat)
	}

	return nil
}

// sanitizeFileName replaces characters that aren't valid in file names on common filesystems
// (e.g. "/" or "\\", which would otherwise be misread as path separators).
func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}

func writeStringToFile(content string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create chapter file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write chapter file: %v", err)
	}

	return nil
}
