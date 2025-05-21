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

func AnalyzeHtmlUrlAndLink(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("Analyzing Html URL and Link function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("URL and Link analyzer succesfully completed in %v", time.Since(start))
	}(startTime)

	doc, err := html.Parse(strings.NewReader(wc.Content))
	if err != nil {
		log.Println("Failed to decode HTML while analyze URL and Link:", err)
		return &response.ErrorResponse{
			Message:  "Failed to decode HTML while analyze URL and Link",
			ErrorMsg: err.Error(),
			Code:     http.StatusBadRequest,
		}
	}

	// Create LinkAnalyzeData to collect links
	var data LinkAnalyzeData

	// Extract URLs into data.Links
	ExtractURL(doc, res.BasePath, &data)
	log.Printf("Total %d links in web content", len(data.Links))
	CheckLinkIsAccessible(&data, res)
	return nil
}

// extractURL extracts URLs from the HTML document.
func ExtractURL(n *html.Node, base string, data *LinkAnalyzeData) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == constant.H_REF || a.Key == constant.SRC {
				// Resolve relative URLs.
				absUrl := UrlResolver(a.Val, base)
				// Only add if not empty after resolving
				if absUrl != constant.EMPTY {
					// collecting links
					data.Links = append(data.Links, absUrl)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ExtractURL(c, base, data)
	}
}

// UrlResolver resolves a URL relative to a base URL.
func UrlResolver(urlStr, base string) string {
	// check url first character is #
	if strings.HasPrefix(urlStr, constant.HASH_CODE) {
		return constant.EMPTY
	}

	// chec url parser has error
	urlPaserCheck, err := url.Parse(urlStr)
	if err != nil {
		return constant.EMPTY
	}

	// check url schema is not empty
	if urlPaserCheck.IsAbs() {
		return urlStr
	}

	baseUrl, err := url.Parse(base)
	if err != nil {
		// invalid base URL
		return constant.EMPTY
	}

	// check provided url is relative or absolute
	resolved, err := baseUrl.Parse(urlStr)
	if err != nil {
		// invalid relative URL
		return constant.EMPTY
	}
	return resolved.String()
}

func CheckLinkIsAccessible(data *LinkAnalyzeData, res *response.SuccessResponse) {
	log.Println("CheckLinkIsAccessible function is executed...")

	// make channel to collect urls
	urlChan := make(chan response.Url)

	// collector of goroutine
	go func() {
		for urlData := range urlChan {
			res.Urls = append(res.Urls, urlData)
		}
	}()

	var wg sync.WaitGroup

	for _, link := range data.Links {
		link := link
		wg.Add(1)
		go func() {
			defer wg.Done()
			urls := response.Url{
				Url: link,
			}

			// check link type
			if strings.Contains(link, res.BasePath) {
				urls.Type = constant.INTERNAL
			} else {
				urls.Type = constant.EXTERNAL
			}

			// check the link accessability
			startTime := time.Now()
			response, err := configs.GetConfig().Client.Get(link)
			if err == nil {
				if response.StatusCode == 200 {
					urls.Accessible = true
				}
				urls.Status = response.StatusCode
			} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				log.Println("Timeout detected using net.Error.Timeout()")
				urls.Status = 408
			}
			urls.UrlExecutionTime = time.Since(startTime).Milliseconds()

			// return date to add to the channel
			urlChan <- urls
		}()
	}

	// wait until all channel complete
	wg.Wait()
	close(urlChan)
}
