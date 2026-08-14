package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "wattpad-parser: fetches and parses a Wattpad story page.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "  %s -url <story-url>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	url := flag.String("url", "", "wattpad story base page URL (e.g https://www.wattpad.com/story/138202651-duplicity-h-s) (required)")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "error: -url flag is required")
		os.Exit(1)
	}

	glossaryText, err := fetchRawHtml(*url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to fetch glossary:", err)
		os.Exit(1)
	}

	fmt.Println(glossaryText)
}
