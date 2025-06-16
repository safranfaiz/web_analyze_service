package analyze

import (
	"api/response"
	"log"
	"strings"
	"time"
)

// HtmlVersionAnalyzer implements the Analyzer interface for HTML versions.
type HtmlVersionAnalyzer struct {
	types map[string]string
}

// NewHtmlVersionAnalyzer creates a new HtmlVersionAnalyzer.
func NewHtmlVersionAnalyzer() *HtmlVersionAnalyzer {
	return &HtmlVersionAnalyzer{
		types: map[string]string{
			"HTML 5":                 `<!DOCTYPE html>`,
			"HTML 4.01 Strict":       `"-//W3C//DTD HTML 4.01//EN"`,
			"HTML 4.01 Transitional": `"-//W3C//DTD HTML 4.01 Transitional//EN"`,
			"HTML 4.01 Frameset":     `"-//W3C//DTD HTML 4.01 Frameset//EN"`,
			"XHTML 1.0 Strict":       `"-//W3C//DTD XHTML 1.0 Strict//EN"`,
			"XHTML 1.0 Transitional": `"-//W3C//DTD XHTML 1.0 Transitional//EN"`,
			"XHTML 1.0 Frameset":     `"-//W3C//DTD XHTML 1.0 Frameset//EN"`,
			"XHTML 1.1":              `"-//W3C//DTD XHTML 1.1//EN"`,
		},
	}
}

// Analyze performs HTML version analysis on the web content.
func (a *HtmlVersionAnalyzer) Analyze(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("Analyzing HTML version function is started...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("HtmlVersionAnalyzer.Analyze completed. Time taken : %d Microseconds", time.Since(startTime).Microseconds())
	}(startTime)

	htmlContent := strings.ToLower(wc.Content)
	for key, val := range a.types {
		if strings.Contains(htmlContent, strings.ToLower(val)) {
			res.HtmlVersion = key
			break
		}
	}
	return nil
}
