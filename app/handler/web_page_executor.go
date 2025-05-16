package handler

import (
	"api/analyze"
	"api/response"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WebPageExecutorHandler(c *gin.Context) {
	startTime := time.Now()
	url := c.Query("url")
	log.Println("executed web page url :", url)
	if url == "" {
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("URL is not exist", nil))
		return
	}

	resp, err := http.Get(url)
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

	resTime := time.Since(startTime).Milliseconds()
	res := &response.SuccessResponse {
		WebPageExtractTime: resTime,
	}
	log.Printf("Web page analysis success with time: %d ms", resTime)

	wc := &response.WebContent{
		Content: string(body),
	}

	analyze.AnalyzeHtmlVersion(wc, res)
	analyze.AnalyzeHtmlTitle(wc, res)
	analyze.AnalyzeHtmlLoginForm(wc, res)
	analyze.AnalyzeHtmlHeading(wc, res)

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}