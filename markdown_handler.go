package main

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func convertHtmlToMarkdown(chapterRawHtml string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(chapterRawHtml)
	if err != nil {
		return "", err
	}

	return markdown, nil
}
