package main

import "fmt"

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatHtml     OutputFormat = "html"
	FormatText     OutputFormat = "text"
)

type RunConfig struct {
	url       string
	outputDir string
	format    OutputFormat
}

// duDE, this implicit interface satisfaction is the bee's knees. It allows us to use the RunConfig struct as a flag.Value, enabling us to define a custom flag for output format.  golang cooked, no notes
func (f *OutputFormat) String() string {
	return string(*f)
}

func (f *OutputFormat) Set(value string) error {
	switch value {
	case "markdown", "html", "text":
		*f = OutputFormat(value)
		return nil
	default:
		return fmt.Errorf("invalid output format: %s. Valid options are: markdown, html, text", value)
	}
}
