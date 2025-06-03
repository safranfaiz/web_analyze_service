package test

import (
	"api/analyze"
	"api/response"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestHtmlHeadingAnalyzer_Analyze_SingleH1(t *testing.T) {
	htmlContent := `<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Main Heading</h1></body></html>`
	wc := &response.WebContent{Content: htmlContent}
	res := &response.SuccessResponse{}
	analyzer := analyze.NewHtmlHeadingAnalyzer()

	err := analyzer.Analyze(wc, res)

	assert.Nil(t, err)
	assert.Len(t, res.Headings, 1)
	if len(res.Headings) == 1 {
		assert.Equal(t, "h1", res.Headings[0].Tag)
		assert.Equal(t, "Main Heading", res.Headings[0].Text)
	}
}

func TestHtmlHeadingAnalyzer_Analyze_MultipleHeadings(t *testing.T) {
	htmlContent := `<html><body><h1>Title 1</h1><p>Some text</p><h2>Subtitle 2</h2></body></html>`
	wc := &response.WebContent{Content: htmlContent}
	res := &response.SuccessResponse{}
	analyzer := analyze.NewHtmlHeadingAnalyzer()

	err := analyzer.Analyze(wc, res)

	assert.Nil(t, err)
	assert.Len(t, res.Headings, 2)
	if len(res.Headings) == 2 {
		foundH1 := false
		foundH2 := false
		for _, h := range res.Headings {
			if h.Tag == "h1" && h.Text == "Title 1" {
				foundH1 = true
			}
			if h.Tag == "h2" && h.Text == "Subtitle 2" {
				foundH2 = true
			}
		}
		assert.True(t, foundH1, "H1 heading not found or incorrect")
		assert.True(t, foundH2, "H2 heading not found or incorrect")
	}
}

func TestHtmlHeadingAnalyzer_Analyze_NoHeadings(t *testing.T) {
	htmlContent := `<html><body><p>Just a paragraph.</p></body></html>`
	wc := &response.WebContent{Content: htmlContent}
	res := &response.SuccessResponse{}
	analyzer := analyze.NewHtmlHeadingAnalyzer()

	err := analyzer.Analyze(wc, res)

	assert.Nil(t, err)
	assert.Len(t, res.Headings, 0)
}

func TestHtmlHeadingAnalyzer_Analyze_NestedContent(t *testing.T) {
	htmlContent := `<html><body><h1>Heading with <strong>bold</strong> text</h1></body></html>`
	wc := &response.WebContent{Content: htmlContent}
	res := &response.SuccessResponse{}
	analyzer := analyze.NewHtmlHeadingAnalyzer()

	err := analyzer.Analyze(wc, res)

	assert.Nil(t, err)
	assert.Len(t, res.Headings, 1)
	if len(res.Headings) == 1 {
		assert.Equal(t, "h1", res.Headings[0].Tag)
		assert.Equal(t, "Heading with", strings.TrimSpace(res.Headings[0].Text))
	}
}

func TestHtmlHeadingAnalyzer_Analyze_DeeplyNestedText(t *testing.T) {
	htmlContent := `<html><body><h1><span><strong>Deep Text</strong></span></h1></body></html>`
	wc := &response.WebContent{Content: htmlContent}
	res := &response.SuccessResponse{}
	analyzer := analyze.NewHtmlHeadingAnalyzer()

	err := analyzer.Analyze(wc, res)

	assert.Nil(t, err)
	assert.Len(t, res.Headings, 1)
	if len(res.Headings) == 1 {
		assert.Equal(t, "h1", res.Headings[0].Tag)
		assert.Equal(t, "span", strings.TrimSpace(res.Headings[0].Text))
	}
}

func TestExactRegexPatternAndToken(t *testing.T) {
	regex := regexp.MustCompile(`h[1-6]`)

	// Test case 1: Match h1
	tokenizerH1 := html.NewTokenizer(strings.NewReader(`<h1>Text</h1>`))
	tokenizerH1.Next() // Move to StartTagToken
	tokenH1, matchH1 := analyze.ExactRegexPatternAndToken(tokenizerH1, regex)
	assert.True(t, matchH1)
	assert.Equal(t, "h1", tokenH1.Data)

	// Test case 2: Match h3
	tokenizerH3 := html.NewTokenizer(strings.NewReader(`<h3 class="foo">Text</h3>`))
	tokenizerH3.Next()
	tokenH3, matchH3 := analyze.ExactRegexPatternAndToken(tokenizerH3, regex)
	assert.True(t, matchH3)
	assert.Equal(t, "h3", tokenH3.Data)

	// Test case 3: No match (p tag)
	tokenizerP := html.NewTokenizer(strings.NewReader(`<p>Text</p>`))
	tokenizerP.Next()
	tokenP, matchP := analyze.ExactRegexPatternAndToken(tokenizerP, regex)
	assert.False(t, matchP)
	assert.Equal(t, "p", tokenP.Data)

	// Test case 4: No match (div tag)
	tokenizerDiv := html.NewTokenizer(strings.NewReader(`<div>Text</div>`))
	tokenizerDiv.Next()
	tokenDiv, matchDiv := analyze.ExactRegexPatternAndToken(tokenizerDiv, regex)
	assert.False(t, matchDiv)
	assert.Equal(t, "div", tokenDiv.Data)
}

func TestSetHeadingDataToResponse(t *testing.T) {
	res := &response.SuccessResponse{}
	textToken := html.Token{Type: html.TextToken, Data: "Heading Text"}
	tagToken := html.Token{Type: html.StartTagToken, Data: "h2"}

	analyze.SetHeadingDataToResponse(textToken, res, tagToken)

	assert.Len(t, res.Headings, 1)
	assert.Equal(t, "h2", res.Headings[0].Tag)
	assert.Equal(t, "Heading Text", res.Headings[0].Text)

	// Add another one
	textToken2 := html.Token{Type: html.TextToken, Data: "Another Heading"}
	tagToken2 := html.Token{Type: html.StartTagToken, Data: "h3"}
	analyze.SetHeadingDataToResponse(textToken2, res, tagToken2)

	assert.Len(t, res.Headings, 2)
	assert.Equal(t, "h3", res.Headings[1].Tag)
	assert.Equal(t, "Another Heading", res.Headings[1].Text)
}
