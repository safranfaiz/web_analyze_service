package test

import (
	"api/app/handler"
	"api/configs"
	"api/constant"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const PATH = "/api/v1/analyze?url="

func TestMain(m *testing.M) {
	// Set TEST_ENV to true for all tests in this package
	// This ensures that configs.loadConfig() picks up test.env
	// if GetConfig() is called.
	os.Setenv(constant.TEST_ENV, "true")
	exitVal := m.Run()             // Run all tests in the package
	os.Unsetenv(constant.TEST_ENV) // Clean up
	os.Exit(exitVal)
}

// Mocking io.ReadCloser for testing handleResponseBodyRead
type mockReadCloser struct {
	reader io.Reader
	closed bool
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.closed {
		return 0, errors.New("read on closed reader")
	}
	return m.reader.Read(p)
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

type mockFailingReadCloser struct{}

func (m *mockFailingReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}
func (m *mockFailingReadCloser) Close() error { return nil }

// --- Existing Tests ---
func TestWebPageExecutorHandlerUrlIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, PATH, nil)

	handler.WebPageExecutorHandler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"response":{"message":"URL is not exist","errorMsg":null,"statusCode":400}}`, w.Body.String())
}

func TestWebPageExecutorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	// Mock the target server that CallWebUrl will hit
	mockTargetServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "<html><head><title>Mocked Page</title></head><body>Hello</body></html>")
	}))
	defer mockTargetServer.Close()

	// Temporarily replace the global HTTP client with the test server's client
	originalClient := configs.GetConfig().Client
	configs.GetConfig().Client = mockTargetServer.Client()
	defer func() { configs.GetConfig().Client = originalClient }()

	c, _ := gin.CreateTestContext(w)
	// Use the mockTargetServer.URL for the request
	c.Request = httptest.NewRequest(http.MethodGet, PATH+mockTargetServer.URL, nil)

	handler.WebPageExecutorHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- New Tests for Helper Functions ---

func TestValidateWebUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test case 1: Valid URL
	wValid := httptest.NewRecorder()
	cValid, _ := gin.CreateTestContext(wValid)
	validURL := "http://example.com"
	resValid, notValid := handler.ValidateWebUrl(validURL, cValid)
	assert.False(t, notValid)
	assert.NotNil(t, resValid)
	assert.Equal(t, validURL, resValid.ExecutedUrl)
	assert.Equal(t, "http://example.com", resValid.BasePath) // Check BasePath extraction

	// Test case 2: Empty URL
	wEmpty := httptest.NewRecorder()
	cEmpty, _ := gin.CreateTestContext(wEmpty)
	resEmpty, isEmpty := handler.ValidateWebUrl("", cEmpty)
	assert.True(t, isEmpty)
	assert.Nil(t, resEmpty)
	assert.Equal(t, http.StatusBadRequest, wEmpty.Code)
	assert.Contains(t, wEmpty.Body.String(), "URL is not exist")

	// Test case 3: URL with path and query
	wPathQuery := httptest.NewRecorder()
	cPathQuery, _ := gin.CreateTestContext(wPathQuery)
	pathQueryURL := "https://sub.example.com/path/page?query=true"
	resPathQuery, notPathQueryValid := handler.ValidateWebUrl(pathQueryURL, cPathQuery)
	assert.False(t, notPathQueryValid)
	assert.NotNil(t, resPathQuery)
	assert.Equal(t, pathQueryURL, resPathQuery.ExecutedUrl)
	assert.Equal(t, "https://sub.example.com", resPathQuery.BasePath)
}

func TestCallWebUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock server to simulate web responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Success!")
		} else if r.URL.Path == "/error" {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()

	originalClient := configs.GetConfig().Client
	configs.GetConfig().Client = mockServer.Client() // Use mock server's client
	defer func() { configs.GetConfig().Client = originalClient }()

	// Test case 1: Successful call
	wSuccess := httptest.NewRecorder()
	cSuccess, _ := gin.CreateTestContext(wSuccess)
	respSuccess, errSuccess := handler.CallWebUrl(mockServer.URL+"/success", cSuccess)
	assert.False(t, errSuccess)
	assert.NotNil(t, respSuccess)
	if respSuccess != nil {
		assert.Equal(t, http.StatusOK, respSuccess.StatusCode)
		bodyBytes, _ := io.ReadAll(respSuccess.Body)
		respSuccess.Body.Close()
		assert.Contains(t, string(bodyBytes), "Success!")
	}

	// Test case 2: Call resulting in server error (client perspective, not a client.Get error)
	// This tests if CallWebUrl correctly returns the response even if it's a 500, etc.
	wServerError := httptest.NewRecorder()
	cServerError, _ := gin.CreateTestContext(wServerError)
	respServerError, errServerError := handler.CallWebUrl(mockServer.URL+"/error", cServerError)
	assert.False(t, errServerError) // client.Get itself didn't fail
	assert.NotNil(t, respServerError)
	if respServerError != nil {
		assert.Equal(t, http.StatusInternalServerError, respServerError.StatusCode)
		respServerError.Body.Close()
	}

	// Test case 3: Call to a non-existent URL (should return an error)
	wClientError := httptest.NewRecorder()
	cClientError, _ := gin.CreateTestContext(wClientError)
	// Restore original client temporarily to make it fail for a non-mocked URL
	configs.GetConfig().Client = http.DefaultClient
	_, errClient := handler.CallWebUrl("http://nonexistentdomain123abc.invalid", cClientError)
	configs.GetConfig().Client = mockServer.Client() // Put back mock server client for other tests if any

	assert.True(t, errClient)
	assert.Equal(t, http.StatusBadRequest, wClientError.Code)
	assert.Contains(t, wClientError.Body.String(), "Error occurred while call web page url")
}

func TestHandleResponseBodyRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test case 1: Successful read
	wSuccess := httptest.NewRecorder()
	cSuccess, _ := gin.CreateTestContext(wSuccess)
	successBodyContent := "Hello, world!"
	mockRespSuccess := &http.Response{
		Body: &mockReadCloser{reader: strings.NewReader(successBodyContent)},
	}
	bodyBytes, errBool := handler.HandleResponseBodyRead(mockRespSuccess, cSuccess)
	assert.False(t, errBool)
	assert.Equal(t, successBodyContent, string(bodyBytes))

	// Test case 2: Error during io.ReadAll
	wError := httptest.NewRecorder()
	cError, _ := gin.CreateTestContext(wError)
	mockRespError := &http.Response{
		Body: &mockFailingReadCloser{},
	}
	_, errBoolError := handler.HandleResponseBodyRead(mockRespError, cError)
	assert.True(t, errBoolError)
	assert.Equal(t, http.StatusBadRequest, wError.Code)
	assert.Contains(t, wError.Body.String(), "Error occurred while reading body")
	assert.Contains(t, wError.Body.String(), "simulated read error")
}
