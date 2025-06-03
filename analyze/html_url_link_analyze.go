package analyze

import (
	"api/configs"
	"api/constant"
	"api/response"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type LinkAnalyzeData struct {
	Links []string
}

// HtmlUrlLinkAnalyzer implements the Analyzer interface for HTML URLs and links.
type HtmlUrlLinkAnalyzer struct{}

// NewHtmlUrlLinkAnalyzer creates a new HtmlUrlLinkAnalyzer.
func NewHtmlUrlLinkAnalyzer() *HtmlUrlLinkAnalyzer {
	return &HtmlUrlLinkAnalyzer{}
}

// Analyze parses HTML, extracts all URLs, and checks their accessibility.
func (a *HtmlUrlLinkAnalyzer) Analyze(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("🔍 Starting analysis of HTML URLs and links...")
	startTime := time.Now()

	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(wc.Content))
	if err != nil {
		log.Println("❌ Failed to parse HTML content:", err)
		return &response.ErrorResponse{
			Message:  "Failed to decode HTML while analyzing URLs",
			ErrorMsg: err.Error(),
			Code:     http.StatusBadRequest,
		}
	}

	// Extract links
	var data LinkAnalyzeData
	extractLinks(doc, res.BasePath, &data)
	log.Printf("📎 Found %d links", len(data.Links))

	// Check accessibility
	checkLinkAccessibility(data.Links, res)
	log.Printf("✅ Completed URL and Link analysis in %d ms", time.Since(startTime).Milliseconds())
	return nil
}

// extractLinks recursively traverses the DOM tree and collects href/src links.
func extractLinks(n *html.Node, base string, data *LinkAnalyzeData) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == constant.H_REF || attr.Key == constant.SRC {
				if absURL := resolveURL(attr.Val, base); absURL != constant.EMPTY {
					data.Links = append(data.Links, absURL)
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		extractLinks(child, base, data)
	}
}

// resolveURL resolves relative URLs against the base and filters anchors or invalid URLs.
func resolveURL(rawURL, base string) string {
	if strings.HasPrefix(rawURL, constant.HASH_CODE) {
		return constant.EMPTY
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return constant.EMPTY
	}

	if parsedURL.IsAbs() {
		return parsedURL.String()
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return constant.EMPTY
	}

	resolvedURL := baseURL.ResolveReference(parsedURL)
	return resolvedURL.String()
}

// checkLinkAccessibility checks which links are accessible and classifies them as internal/external.
func checkLinkAccessibility(links []string, res *response.SuccessResponse) {
	log.Println("🌐 Checking link accessibility...")

	urlChan := make(chan response.Url)
	var wg sync.WaitGroup

	// Collector goroutine
	go func() {
		for urlData := range urlChan {
			res.Urls = append(res.Urls, urlData)
		}
	}()

	client := configs.GetConfig().Client

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			checkSingleURL(link, res.BasePath, client, urlChan)
		}(link)
	}

	wg.Wait()
	close(urlChan)
}

// checkSingleURL checks the accessibility of a single URL and sends the result through a channel.
func checkSingleURL(link, basePath string, client *http.Client, urlChan chan<- response.Url) {
	result := response.Url{
		Url:  link,
		Type: classifyLinkType(link, basePath),
	}

	start := time.Now()
	resp, err := client.Get(link)
	result.UrlExecutionTime = time.Since(start).Milliseconds()

	if err == nil {
		result.Status = resp.StatusCode
		result.Accessible = resp.StatusCode == http.StatusOK
	} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		log.Printf("⏰ Timeout accessing: %s", link)
		result.Status = http.StatusRequestTimeout
	} else {
		log.Printf("⚠️ Failed accessing: %s | Error: %v", link, err)
	}

	urlChan <- result
}

// classifyLinkType determines whether the link is internal or external.
func classifyLinkType(link, base string) string {
	if strings.Contains(link, base) {
		return constant.INTERNAL
	}
	return constant.EXTERNAL
}
