package analyze

import (
	"api/constant"
	"api/response"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// AnalyzeHtmlTitle is responsible for set the response to HTML Title tag text of given URL
func AnalyzeHtmlTitle(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("Analyzing HTML title function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("Title analyzer succesfully completed in %d ms", time.Since(start).Microseconds())
	}(startTime)

	reader, err := charset.NewReader(strings.NewReader(wc.Content), constant.EMPTY)
	if err != nil {
		log.Println("Failed to decode HTML:", err)
		return &response.ErrorResponse{
			Message:  "Failed to decode HTML while Analyze Html Title",
			ErrorMsg: err.Error(),
			Code:     http.StatusBadRequest,
		}
	}
	metaData := html.NewTokenizer(reader)

	for {
		token := metaData.Next()
		switch token {
		case html.ErrorToken:
			log.Printf("HTML content having error and title analyzer stop in %d ms", time.Since(startTime))
			return &response.ErrorResponse{
				Message: "HTML content having error while Analyze Html Title",
				Code:    http.StatusBadRequest,
			}
		case html.StartTagToken:
			t := metaData.Token()
			if t.Data == constant.TAG_TITLE {
				if metaData.Next() == html.TextToken {
					res.Title = metaData.Token().Data
					return nil
				}
			}
		}
	}
}
