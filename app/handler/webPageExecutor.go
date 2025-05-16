package handler

import (
	"api/response"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebPageExecutorHandler(c *gin.Context) {
	url := c.Query("url")
	log.Println("executed web page url :", url)
	if url == "" {
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("No URL for execute ", nil))
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error occurred while call web page url : ", err)
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("Error occurred while call web page url ", err))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occurred : ", err)
		c.IndentedJSON(http.StatusBadRequest, response.ErrorResponseMsg("Error occurred while reading body ", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": string(body),
	})
}