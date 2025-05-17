package handler

import (
	"api/analyze"
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
	if link == constant.EMPTY {
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("URL is not exist", nil))
		return
	}

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal("Error occurred while call web page url", err)
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("Error occurred while call web page url", err))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error occurred while reading response body", err)
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("Error occurred while reading body", err))
		return
	}

	parsedURL, _ := url.Parse(link)
	baseUrl := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	resTime := time.Since(startTime).Milliseconds()
	res := &response.SuccessResponse{
		WebPageExtractTime: resTime,
		ExecutedUrl:        link,
		BasePath:           baseUrl,
	}
	log.Printf("Web page analysis success with time: %d ms", resTime)

	wc := &response.WebContent{
		Content: string(body),
	}

	analyze.AnalyzeHtmlVersion(wc, res)
	analyze.AnalyzeHtmlTitle(wc, res)
	analyze.AnalyzeHtmlLoginForm(wc, res)
	analyze.AnalyzeHtmlHeading(wc, res)
	analyze.AnalyzeHtmlUrlAndLink(wc, res)

	appExecuteTotalTime := time.Since(startTime).Milliseconds()
	res.AppExecuteTotalTime = appExecuteTotalTime
	c.JSON(http.StatusOK, gin.H{
		"response": res,
	})
}
