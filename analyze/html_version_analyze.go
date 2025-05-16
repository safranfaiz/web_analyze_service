package analyze

import (
	"api/response"
	"log"
	"strings"
	"time"
)

type htmlVersionAnalyzer struct {
	types map[string]string
}

func HtmlVersions() htmlVersionAnalyzer {
	return htmlVersionAnalyzer{
		types: map[string]string{
			"HTML 5":                  `<!DOCTYPE html>`,
			"HTML 4.01 Strict":        `"-//W3C//DTD HTML 4.01//EN"`,
			"HTML 4.01 Transitional":  `"-//W3C//DTD HTML 4.01 Transitional//EN"`,
			"HTML 4.01 Frameset":      `"-//W3C//DTD HTML 4.01 Frameset//EN"`,
			"XHTML 1.0 Strict":        `"-//W3C//DTD XHTML 1.0 Strict//EN"`,
			"XHTML 1.0 Transitional":  `"-//W3C//DTD XHTML 1.0 Transitional//EN"`,
			"XHTML 1.0 Frameset":      `"-//W3C//DTD XHTML 1.0 Frameset//EN"`,
			"XHTML 1.1":               `"-//W3C//DTD XHTML 1.1//EN"`,
		},
	}
}

// AnalyzeHtmlVersion is responsible for set the response to HTML verion of given URL
func AnalyzeHtmlVersion(wc *response.WebContent, res *response.SuccessResponse) {
	log.Println("Analyzing HTML version function is started...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("Analyzing HTML version function completed. Time taken : %v ms", time.Since(start).Milliseconds())
	}(startTime)

	htmlContent := strings.ToLower(wc.Content)
	for key, val := range HtmlVersions().types {
		if strings.Contains(htmlContent, strings.ToLower(val)) {
			res.HtmlVersion = key
			break
		}
	}
}
