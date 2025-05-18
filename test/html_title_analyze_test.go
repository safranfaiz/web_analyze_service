package test

import (
	"api/analyze"
	"api/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeHtmlTitleSuccess(t *testing.T) {
	htmlContent := "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n\t<meta charset='utf-8'>\n\t<meta name='viewport' content='width=device-width,initial-scale=1'>\n\n\t<title>Test Web Page Analyzer</title>\n</head>\n\n<body>\n</body>\n</html>"
	wc := response.WebContent{
		Content: htmlContent,
	}
	res := response.SuccessResponse{}
	analyze.AnalyzeHtmlTitle(&wc, &res)
	assert.Equal(t, "Test Web Page Analyzer", res.Title, "Web Page Title Testing...")
}

func TestAnalyzeHtmlTitlFailInReader(t *testing.T) {
	htmlContent := ""
	wc := response.WebContent{
		Content: htmlContent,
	}
	res := response.SuccessResponse{}
	err := analyze.AnalyzeHtmlTitle(&wc, &res)
	assert.Equal(t, "Failed to decode HTML while Analyze Html Title", err.Message, "Web Page Title Testing...")
}

func TestAnalyzeHtmlTitlFailInErrorToken(t *testing.T) {
	htmlContent := "<html><body><p>Hello"
	wc := response.WebContent{
		Content: htmlContent,
	}
	res := response.SuccessResponse{}
	err := analyze.AnalyzeHtmlTitle(&wc, &res)
	assert.Equal(t, "HTML content having error while Analyze Html Title", err.Message, "Web Page Title Testing...")
}
