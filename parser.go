package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type remixContext struct {
	State struct {
		LoaderData struct {
			Story struct {
				Story struct {
					Parts []struct {
						ID    int    `json:"id"`
						Title string `json:"title"`
						URL   string `json:"url"`
					} `json:"parts"`
				} `json:"story"`
			} `json:"routes/story.$storyid"`
		} `json:"loaderData"`
	} `json:"state"`
}

type chapterRecord struct {
	title string
	url   string
}

func extractChapterURLs(rootPageRawHtml string) ([]chapterRecord, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rootPageRawHtml))
	if err != nil {
		return nil, err
	}

	var scriptText string
	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := s.Text()
		if strings.Contains(text, "window.__remixContext") {
			scriptText = text
			return false // stop iterating
		}
		return true
	})

	if scriptText == "" {
		return nil, fmt.Errorf("could not find __remixContext script")
	}

	// isolate the JSON assigned to window.__remixContext
	const prefix = "window.__remixContext = "
	_, after, ok := strings.Cut(scriptText, prefix)
	if !ok {
		return nil, fmt.Errorf("could not find __remixContext assignment")
	}
	jsonStr := strings.TrimSpace(after)
	jsonStr = strings.TrimSuffix(jsonStr, ";")

	var ctx remixContext
	if err := json.Unmarshal([]byte(jsonStr), &ctx); err != nil {
		return nil, err
	}

	urls := make([]chapterRecord, 0, len(ctx.State.LoaderData.Story.Story.Parts))
	for _, p := range ctx.State.LoaderData.Story.Story.Parts {
		urls = append(urls, chapterRecord{
			title: p.Title,
			url:   p.URL,
		})
	}
	return urls, nil
}
