package handler

import (
	"api/analyze"
	"api/configs"
	"api/constant"
	"api/response"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func WebPageExecutorHandler(c *gin.Context) {
	startTime := time.Now()
	link := c.Query(constant.URL)
	log.Println("executed web page url :", link)

	res, notValid := ValidateWebUrl(link, c)
	if notValid {
		return
	}

	resp, webUrlError := CallWebUrl(link, c)
	if webUrlError {
		return
	}

	body, err := handleResponseBodyRead(resp, c)
	if err {
		return
	}

	resTime := time.Since(startTime).Milliseconds()
	res.WebPageExtractTime = resTime
	log.Printf("Web page analysis success with time: %d ms", resTime)

	wc := &response.WebContent{
		Content: string(body),
	}

	// execute analyzers for collecting meta data of web page
	analyze.AnalyzeHtmlVersion(wc, res)
	analyze.AnalyzeHtmlTitle(wc, res)
	analyze.AnalyzeHtmlLoginForm(wc, res)
	analyze.AnalyzeHtmlHeading(wc, res)
	analyze.AnalyzeHtmlUrlAndLink(wc, res)

	appExecuteTotalTime := time.Since(startTime).Milliseconds()
	res.AppExecuteTotalTime = appExecuteTotalTime
	c.JSON(http.StatusOK, gin.H{
		constant.RESPONSE: res,
	})
}

func handleResponseBodyRead(resp *http.Response, c *gin.Context) ([]byte, bool) {
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
