package analyze

import (
	"api/response"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const headingHTMLTagRegex = `h[1-6]`

func AnalyzeHtmlHeading(wc *response.WebContent, res *response.SuccessResponse) {
	log.Println("Analyzing Login form function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("Login form analyzer succesfully completed in %v", time.Since(start))
	}(startTime)

	metaData := html.NewTokenizer(strings.NewReader(wc.Content))

	for {
		tagType := metaData.Next()
		switch tagType {
		case html.StartTagToken:
			token := metaData.Token()
			if isHeadingTag(token.Data) {
				textContent := extractTextContent(metaData)
				log.Println("Tag: ",token.Data, " Level: ",textContent)
				// set the heading tag and that content to response
				
			}
		case html.ErrorToken:
			err := metaData.Err()
			if err != nil {
				log.Fatalf("HTML tokenizer error: %v", err)
			}
			return
		}
	}
}

// isHeadingTag function encapsulates the logic for checking if a tag name matches the heading tag regex
func isHeadingTag(tagName string) bool {
	regex, err := regexp.Compile(headingHTMLTagRegex)
	if err != nil {
		log.Println("failed to compile heading regex: " + err.Error())
	}
	return regex.MatchString(tagName)
}

// extractTextContent function uses a depth counter.
// when a StartTagToken is encountered inside a heading tag, the depth is incremented.
// When an EndTagToken is encountered, the depth is decremented.
// The loop breaks when the depth becomes less than 0, indicating that the 
// closing tag of the initial heading tag has been reached.
func extractTextContent(tokenizer *html.Tokenizer) string {
	var textContent strings.Builder
	depth := 0

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.TextToken:
			if depth == 0 {
				textContent.WriteString(tokenizer.Token().Data)
			}
		case html.StartTagToken:
			depth++
		case html.EndTagToken:
			depth--
			if depth < 0 {
				return strings.TrimSpace(textContent.String())
			}
		case html.ErrorToken:
			return strings.TrimSpace(textContent.String())
		}
	}
}