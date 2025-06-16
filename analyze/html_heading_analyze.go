package analyze

import (
	"api/response"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const headingHTMLTagRegex = `h[1-6]`

// HtmlHeadingAnalyzer implements the Analyzer interface for HTML headings.
type HtmlHeadingAnalyzer struct{}

// NewHtmlHeadingAnalyzer creates a new HtmlHeadingAnalyzer.
func NewHtmlHeadingAnalyzer() *HtmlHeadingAnalyzer {
	return &HtmlHeadingAnalyzer{}
}

// Analyze performs heading analysis on the web content.
func (a *HtmlHeadingAnalyzer) Analyze(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("Analyzing HTML Headings function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("HtmlHeadingAnalyzer.Analyze execution completed...")
	}(startTime)

	regex, err := regexp.Compile(headingHTMLTagRegex)
	if err != nil {
		log.Fatal("Error occurred in compiling regex", err)
		return &response.ErrorResponse{
			Message:  "Error occurred in compiling regex",
			ErrorMsg: err.Error(),
			Code:     http.StatusBadRequest,
		}
	}
	metaData := html.NewTokenizer(strings.NewReader(wc.Content))

	return ValidateHeaderTokenHandler(metaData, regex, res, startTime)
}

func ValidateHeaderTokenHandler(metaData *html.Tokenizer, regex *regexp.Regexp, res *response.SuccessResponse, startTime time.Time) *response.ErrorResponse {
OuterLoop:
	for {
		switch metaData.Next() {
		case html.StartTagToken:
			token, match := ExactRegexPatternAndToken(metaData, regex)
			if match {
				// scan next token and return
				metaData.Next()
				tempToken := metaData.Token()
				if tempToken.Type == html.TextToken {
					SetHeadingDataToResponse(tempToken, res, token)
				} else {
					// handle deeper nested or multiline text content
					for {
						switch metaData.Next() {
						case html.TextToken:
							SetHeadingDataToResponse(tempToken, res, token)
							// break out of the outer loop
							break OuterLoop
						case html.ErrorToken:
							err := metaData.Err()
							// EOF mean no more input is available
							if err == io.EOF {
								break OuterLoop
							}
							log.Printf("HTML tokenizer error: %v", err)
							// break out of the outer loop
							break OuterLoop
						}
					}
				}
			}

		case html.ErrorToken:
			err := metaData.Err()
			if err == io.EOF {
				break OuterLoop
			}
			log.Printf("HTML tokenizer error: %v", err)
			break OuterLoop
		}
	}
	log.Printf("Analyzing HTML Headings succesfully completed in %d ms", time.Since(startTime).Milliseconds())
	return nil
}

func SetHeadingDataToResponse(tempToken html.Token, res *response.SuccessResponse, token html.Token) {
	res.Headings = append(res.Headings, response.Heading{
		Tag:  token.Data,
		Text: tempToken.Data,
	})
}

// ExactRegexPatternAndToken is responsible for pattern and token exact
func ExactRegexPatternAndToken(metaData *html.Tokenizer, regex *regexp.Regexp) (html.Token, bool) {
	token := metaData.Token()
	match := regex.Match([]byte(token.Data))
	return token, match
}
