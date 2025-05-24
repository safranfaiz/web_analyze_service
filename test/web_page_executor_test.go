package test

import (
	"api/app/handler"
	"api/constant"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const PATH = "/api/v1/analyze?url="

func TestWebPageExecutorHandlerUrlIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, PATH, nil)

	handler.WebPageExecutorHandler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	//log.Println("Body of response: ", w.Body.String())
	assert.Equal(t, "{\"response\":{\"message\":\"URL is not exist\",\"errorMsg\":null,\"statusCode\":400}}", w.Body.String())
}

func TestWebPageExecutorHandler(t *testing.T) {
	// Set environment variable BEFORE anything else
	os.Setenv(constant.TEST_ENV, "true")
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, PATH+"https://safranfaiz.github.io/shariputhra_maha_vidyalaya_ahangama/", nil)

	handler.WebPageExecutorHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
	os.Unsetenv(constant.TEST_ENV)
}
