package handler

import (
	"api/analyze"
	"api/configs"
	"api/constant"
	"api/response"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func WebPageExecutorHandler(c *gin.Context) {
	startTime := time.Now()
	link := c.Query(constant.URL)
	log.Println("executed web page url :", link)

	res, notValid := ValidateWebUrl(link, c) // Assuming ValidateWebUrl is now exported
	if notValid {
		return
	}

	resp, webUrlError := CallWebUrl(link, c) // Assuming CallWebUrl is now exported
	if webUrlError {
		return
	}

	body, err := HandleResponseBodyRead(resp, c) // Use exported name
	if err {
		return
	}

	resTime := time.Since(startTime).Milliseconds()
	res.WebPageExtractTime = resTime
	log.Printf("Web page analysis success with time: %d ms", resTime)

	wc := &response.WebContent{
		Content: string(body),
	}

	// Create a list of analyzers
	analyzers := []analyze.Analyzer{
		analyze.NewHtmlVersionAnalyzer(),
		analyze.NewHtmlTitleAnalyzer(),
		analyze.NewHtmlLoginFormAnalyzer(),
		analyze.NewHtmlHeadingAnalyzer(),
		analyze.NewHtmlUrlLinkAnalyzer(),
	}

	// Execute analyzers concurrently
	var wg sync.WaitGroup
	errChan := make(chan *response.ErrorResponse, len(analyzers)) // Buffered channel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to free resources

	for _, analyzerInstance := range analyzers {
		wg.Add(1)
		go func(analyzer analyze.Analyzer) {
			defer wg.Done()
			select {
			case <-ctx.Done(): // Check if context was cancelled
				log.Printf("Analysis cancelled for %T due to an error in another analyzer.", analyzer)
				return
			default:
				if analysisErr := analyzer.Analyze(wc, res); analysisErr != nil {
					log.Printf("Error during analysis with %T: %s. Error details: %s", analyzer, analysisErr.Message, analysisErr.ErrorMsg)
					// Try to send error to channel, but don't block if full
					select {
					case errChan <- analysisErr:
						cancel() // Signal other goroutines to stop
					default:
						log.Printf("Error channel full, could not send error from %T", analyzer)
					}
					return
				}
			}
		}(analyzerInstance)
	}

	// Goroutine to close errChan once all analyzers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Wait for the first error or for all to complete
	if firstErr := <-errChan; firstErr != nil {
		// An error occurred in one of the analyzers
		log.Printf("First error received, terminating analysis. Error: %s", firstErr.Message)
		c.JSON(firstErr.Code, gin.H{
			constant.RESPONSE: firstErr,
		})
		return
	}
	// If we reach here, all analyzers completed successfully or were cancelled
	// but no error was sent to errChan before it was closed by wg.Wait().

	appExecuteTotalTime := time.Since(startTime).Milliseconds()
	res.AppExecuteTotalTime = appExecuteTotalTime
	c.JSON(http.StatusOK, gin.H{
		constant.RESPONSE: res,
	})
}

// HandleResponseBodyRead reads the body of an HTTP response.
func HandleResponseBodyRead(resp *http.Response, c *gin.Context) ([]byte, bool) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occurred while reading response body", err)
		c.JSON(http.StatusBadRequest, gin.H{
			constant.RESPONSE: response.ErrorResponseMsg("Error occurred while reading body", err.Error(), http.StatusBadRequest),
		})
		return nil, true
	}
	return body, false
}

// CallWebUrl makes an HTTP GET request to the given link.
func CallWebUrl(link string, c *gin.Context) (*http.Response, bool) {
	resp, err := configs.GetConfig().Client.Get(link)
	if err != nil {
		log.Println("Error occurred while call web page url", err)
		c.JSON(http.StatusBadRequest, gin.H{
			constant.RESPONSE: response.ErrorResponseMsg("Error occurred while call web page url", err.Error(), http.StatusBadRequest),
		})
		return nil, true
	}
	return resp, false
}

// ValidateWebUrl checks if the provided URL string is valid and prepares a SuccessResponse.
func ValidateWebUrl(link string, c *gin.Context) (*response.SuccessResponse, bool) {
	res := &response.SuccessResponse{
		ExecutedUrl: link,
	}
	parsedURL, _ := url.Parse(link)
	baseUrl := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	res.BasePath = baseUrl

	if link == constant.EMPTY {
		log.Println("No URL is exist")
		c.JSON(http.StatusBadRequest, gin.H{
			constant.RESPONSE: response.ErrorResponseMsg("URL is not exist", nil, http.StatusBadRequest),
		})
		return nil, true
	}
	return res, false
}
