package main

import (
	"fmt"
	"os"
)

const (
	// TODO implement this as input to the interface
	TABLE_OF_CONTENTS_URL = "https://www.wattpad.com/story/138202651-duplicity-h-s"
)

func main() {
	// TODO implement "flags" dependency to set command line flags for usage
	// https://gobyexample.com/command-line-flags
	// blah

	glossaryText, err := fetchRawHtml(TABLE_OF_CONTENTS_URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to fetch glossary:", err)
	}

	fmt.Println(glossaryText)
}
