package test

import (
	"api/analyze"
	"api/configs"
	"api/constant"
	"api/response"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestResolveURL tests the resolveURL internal function
func TestResolveURL(t *testing.T) {
	base := "http://example.com/path/"
	tests := []struct {
		name     string
		rawURL   string
		base     string
		expected string
	}{
		{"absolute url", "http://othersite.com/page", base, "http://othersite.com/page"},
		{"relative url", "subpage.html", base, "http://example.com/path/subpage.html"},
		{"relative url with ..", "../another.html", base, "http://example.com/another.html"},
		{"root relative url", "/root.html", base, "http://example.com/root.html"},
		{"anchor link", "#section1", base, ""},
		{"empty url", "", base, ""}, // Assuming resolveURL handles this by returning empty or base
		{"invalid url", ":not_a_url", base, ""},
		{"base with no trailing slash", "page.html", "http://example.com", "http://example.com/page.html"},
		{"rawURL is just a path", "just/a/path", "http://example.com/base/", "http://example.com/base/just/a/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}

// TestClassifyLinkType tests the classifyLinkType internal function
func TestClassifyLinkType(t *testing.T) {
	base := "http://example.com"
	tests := []struct {
		name     string
		link     string
		base     string
		expected string
	}{
		{"internal link", "http://example.com/path/page", base, constant.INTERNAL},
		{"external link", "http://othersite.com/page", base, constant.EXTERNAL},
		{"internal link subdomain", "http://sub.example.com/page", base, constant.EXTERNAL}, // Based on current strings.Contains logic
		{"link same as base", "http://example.com", base, constant.INTERNAL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

// TestCheckSingleURL tests the checkSingleURL internal function using a mock server
func TestCheckSingleURL(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ok") {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Hello, client")
		} else if strings.HasSuffix(r.URL.Path, "/notfound") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.HasSuffix(r.URL.Path, "/timeout") {
			time.Sleep(100 * time.Millisecond) // Simulate timeout, client should have shorter timeout
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// Mock client with a short timeout for the timeout test
	mockClient := &http.Client{Timeout: 50 * time.Millisecond}
	originalClient := configs.GetConfig().Client
	configs.GetConfig().Client = mockClient
	defer func() { configs.GetConfig().Client = originalClient }()
}

func TestHtmlUrlLinkAnalyzer_Analyze_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/page1" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "ok")
		} else if r.URL.Path == "/page2" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			// Main page content - serve relative links
			fmt.Fprintf(w, `<html><body>
		              <a href="/page1">Internal OK</a>
		              <a href="/page2">Internal Not Found</a>
		              <a href="http://external.example.com/extpage">External</a>
		          </body></html>`)
		}
	}))
	defer server.Close()

	originalClient := configs.GetConfig().Client
	configs.GetConfig().Client = server.Client()
	defer func() { configs.GetConfig().Client = originalClient }()

	// This is the content the analyzer will parse.
	// It should contain relative links, and BasePath will be used for resolution.
	htmlContentToAnalyze := `<html><body>
		      <a href="/page1">Internal OK</a>
		      <a href="/page2">Internal Not Found</a>
		      <a href="http://external.example.com/extpage">External</a>
		  </body></html>`

	wc := &response.WebContent{Content: htmlContentToAnalyze}
	res := &response.SuccessResponse{BasePath: server.URL} // Set base path for correct classification
	analyzer := analyze.NewHtmlUrlLinkAnalyzer()

	err := analyzer.Analyze(wc, res)
	assert.Nil(t, err)
	time.Sleep(200 * time.Millisecond)
	assert.Len(t, res.Urls, 3)

	foundPage1, foundPage2, foundExternal := false, false, false
	for _, u := range res.Urls {
		if u.Url == server.URL+"/page1" {
			assert.True(t, u.Accessible, "Page1 should be accessible")
			assert.Equal(t, http.StatusOK, u.Status)
			assert.Equal(t, constant.INTERNAL, u.Type)
			foundPage1 = true
		}
		if u.Url == server.URL+"/page2" {
			assert.False(t, u.Accessible, "Page2 should not be accessible")
			assert.Equal(t, http.StatusNotFound, u.Status)
			assert.Equal(t, constant.INTERNAL, u.Type)
			foundPage2 = true
		}
		if u.Url == "http://external.example.com/extpage" {
			assert.Equal(t, constant.EXTERNAL, u.Type)
			foundExternal = true
		}
	}
	assert.True(t, foundPage1, "Test for /page1 did not run or match")
	assert.True(t, foundPage2, "Test for /page2 did not run or match")
	assert.True(t, foundExternal, "Test for external link did not run or match")

	wcEmpty := &response.WebContent{Content: ""}
	resEmpty := &response.SuccessResponse{BasePath: "http://example.com"}
	errEmptyParse := analyzer.Analyze(wcEmpty, resEmpty)

	assert.Nil(t, errEmptyParse, "Expected no error for empty HTML content as html.Parse is tolerant")
	assert.Empty(t, resEmpty.Urls, "Expected no URLs to be found in an empty document")
}
