package analyze

import (
	"api/constant"
	"api/response"
	"log"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// AnalyzeHtmlTitle is responsible for set the response to HTML Title tag text of given URL
func AnalyzeHtmlTitle(wc *response.WebContent, res *response.SuccessResponse) {
	log.Println("Analyzing HTML title function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("Title analyzer succesfully completed in %v", time.Since(start))
	}(startTime)

	reader, err := charset.NewReader(strings.NewReader(wc.Content), constant.EMPTY)
	if err != nil {
		log.Println("Failed to decode HTML:", err)
		return
	}
	metaData := html.NewTokenizer(reader)

	for {
		token := metaData.Next()
		switch token {
		case html.ErrorToken:
			log.Fatalf("HTML content having error and title analyzer stop in %d ms", time.Since(startTime))
			return
		case html.StartTagToken:
			t := metaData.Token()
			if t.Data == constant.TAG_TITLE {
				if metaData.Next() == html.TextToken {
					res.Title = metaData.Token().Data
					return
				}
			}
		}
	}
}
