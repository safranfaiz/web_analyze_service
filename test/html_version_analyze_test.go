package test

import (
	"api/analyze"
	"api/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeHtmlVersion(t *testing.T) {
	htmlContent := "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n\t<meta charset='utf-8'>\n\t<meta name='viewport' content='width=device-width,initial-scale=1'>\n\n\t<title>Test Web Page Analyzer</title>\n</head>\n\n<body>\n</body>\n</html>"
	wc := response.WebContent{
		Content: htmlContent,
	}
	res := response.SuccessResponse{}
	analyze.AnalyzeHtmlVersion(&wc, &res)
	assert.Equal(t, "HTML 5", res.HtmlVersion, "Web Page Version Testing...")
}
