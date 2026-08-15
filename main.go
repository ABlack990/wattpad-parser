package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	config := RunConfig{
		url:       "",
		outputDir: "./wattpad-output",
		format:    FormatMarkdown,
	}

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "wattpad-parser: Fetches, parses and outputs a Wattpad story as Markdown text.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "  %s -url <story-url> -output-dir <output-directory> -output-format <format>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.StringVar(&config.url, "url", config.url, "wattpad story base page URL (e.g https://www.wattpad.com/story/138202651-duplicity-h-s) (required)")
	flag.StringVar(&config.outputDir, "output-dir", config.outputDir, "directory to save the output files)")
	flag.Var(&config.format, "output-format", "Output format(markdown, html, text)")
	flag.Parse()

	if config.url == "" {
		fmt.Fprintln(os.Stderr, "error: -url flag is required")
		os.Exit(1)
	}

	pageRawHtml, err := fetchRawHtml(config.url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to fetch page:", err)
		os.Exit(1)
	}

	storyMetadata, chapterRecords, err := parseStory(pageRawHtml)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse story:", err)
		os.Exit(1)
	}

	err = os.MkdirAll(config.outputDir, 0775)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create root output directory:", err)
		os.Exit(1)
	}

	err = writeMetadataJsonFile(storyMetadata, config.outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write metadata.json: %v\n", err)
		os.Exit(1)
	}

	for _, chapterRecord := range chapterRecords {
		title := chapterRecord.title
		chapterURL := chapterRecord.url
		pageRawHtml, err := fetchRawHtml(chapterURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch chapter page for %q (%s): %v\n", title, chapterURL, err)
			continue
		}

		chapterRawHtml, err := extractChapterTextHtml(pageRawHtml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to extract chapter text for %q (%s): %v\n", title, chapterURL, err)
			continue
		}

		err = writeChapter(chapterRawHtml, chapterRecord, config.outputDir, config.format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write chapter %q: %v\n", title, err)
			continue
		}
	}
}
