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

type storyPerson struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type storyPublisherLogo struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}

type storyPublisher struct {
	Type             string             `json:"@type"`
	Name             string             `json:"name"`
	LegalName        string             `json:"legalName"`
	FoundingDate     string             `json:"foundingDate"`
	FoundingLocation string             `json:"foundingLocation"`
	AreaServed       string             `json:"areaServed"`
	Logo             storyPublisherLogo `json:"logo"`
	URL              string             `json:"url"`
}

type storyMainEntityOfPage struct {
	Type string `json:"@type"`
	ID   string `json:"@id"`
}

// storyMetadata mirrors the schema.org Article JSON-LD block embedded in a story's root page.
type storyMetadata struct {
	Context              string                `json:"@context"`
	Type                 string                `json:"@type"`
	About                string                `json:"about"`
	Author               storyPerson           `json:"author"`
	AccountablePerson    storyPerson           `json:"accountablePerson"`
	CopyrightHolder      storyPerson           `json:"copyrightHolder"`
	CopyrightYear        string                `json:"copyrightYear"`
	DateCreated          string                `json:"dateCreated"`
	DatePublished        string                `json:"datePublished"`
	DateModified         string                `json:"dateModified"`
	Description          string                `json:"description"`
	CommentCount         int64                 `json:"commentCount"`
	DiscussionURL        string                `json:"discussionUrl"`
	Genre                string                `json:"genre"`
	Headline             string                `json:"headline"`
	InLanguage           string                `json:"inLanguage"`
	InteractionStatistic int64                 `json:"interactionStatistic"`
	IsAccessibleForFree  bool                  `json:"isAccessibleForFree"`
	IsFamilyFriendly     bool                  `json:"isFamilyFriendly"`
	Keywords             string                `json:"keywords"`
	Publisher            storyPublisher        `json:"publisher"`
	PublishingPrinciples string                `json:"publishingPrinciples"`
	ThumbnailURL         string                `json:"thumbnailUrl"`
	TypicalAgeRange      string                `json:"typicalAgeRange"`
	Image                string                `json:"image"`
	Name                 string                `json:"name"`
	URL                  string                `json:"url"`
	MainEntityOfPage     storyMainEntityOfPage `json:"mainEntityOfPage"`
}

func parseStory(rootPageRawHtml string) (storyMetadata, []chapterRecord, error) {
	metadata, err := extractStoryMetadata(rootPageRawHtml)
	if err != nil {
		return metadata, nil, fmt.Errorf("failed to extract story metadata: %v", err)
	}

	chapterRecords, err := extractChapterURLs(rootPageRawHtml)
	if err != nil {
		return metadata, nil, fmt.Errorf("failed to extract chapter URLs: %v", err)
	}

	return metadata, chapterRecords, nil
}

func extractStoryMetadata(rootPageRawHtml string) (storyMetadata, error) {
	var metadata storyMetadata

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rootPageRawHtml))
	if err != nil {
		return metadata, err
	}

	jsonStr := doc.Find(`script[type="application/ld+json"]`).First().Text()
	if jsonStr == "" {
		return metadata, fmt.Errorf("could not find application/ld+json script")
	}

	if err := json.Unmarshal([]byte(jsonStr), &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
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

func extractChapterTextHtml(chapterPageRawHtml string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(chapterPageRawHtml))
	if err != nil {
		return "", err
	}

	pre := doc.Find(".panel-reading pre").First()
	if pre.Length() == 0 {
		return "", fmt.Errorf("could not find chapter text container")
	}

	// strip the injected audio placeholder, not part of the actual chapter content
	pre.Find(".trinityAudioPlaceholder").Remove()

	html, err := pre.Html()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(html), nil
}
