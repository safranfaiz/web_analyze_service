package analyze

import (
	"api/constant"
	"api/response"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// HtmlLoginFormAnalyzer implements the Analyzer interface for HTML login forms.
type HtmlLoginFormAnalyzer struct{}

// NewHtmlLoginFormAnalyzer creates a new HtmlLoginFormAnalyzer.
func NewHtmlLoginFormAnalyzer() *HtmlLoginFormAnalyzer {
	return &HtmlLoginFormAnalyzer{}
}

// Analyze performs login form analysis on the web content.
func (a *HtmlLoginFormAnalyzer) Analyze(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse {
	log.Println("Analyzing Login form function is executed...")
	startTime := time.Now()

	defer func(start time.Time) {
		log.Printf("HtmlLoginFormAnalyzer.Analyze succesfully completed in %v", time.Since(start))
	}(startTime)

	nodes, err := htmlquery.Parse(strings.NewReader(wc.Content))
	if err != nil {
		log.Printf("Parser error while analyze login form and time taken for %v", time.Since(startTime))
		return &response.ErrorResponse{
			Message:  "Parser error while analyze login form",
			ErrorMsg: err.Error(),
			Code:     http.StatusBadRequest,
		}
	}
	forms := htmlquery.Find(nodes, constant.FORM_TAG_EXP)
	ValidateLoginFormHandler(forms, res)
	return nil
}

func ValidateLoginFormHandler(forms []*html.Node, res *response.SuccessResponse) {
	for _, form := range forms {
		var hasUsername, hasPassword, hasSubmit bool

		// check the input type is text or email
		usernameInputs := htmlquery.Find(form, constant.LOGIN_INPUT_VALIDATION)
		if len(usernameInputs) > 0 {
			hasUsername = true
		}

		// check the input type is password
		passwordInputs := htmlquery.Find(form, constant.LOGIN_PASSWORD_VALIDATION)
		if len(passwordInputs) > 0 {
			hasPassword = true
		}

		// check for submit button or input
		submitButtons := htmlquery.Find(form, constant.LOGIN_SUBMIT_BUTTON_VALIDATION)
		if len(submitButtons) > 0 {
			hasSubmit = true
		}

		// login form need user name field, password and submit
		// based on this condition we identiy login form is exist
		if hasUsername && hasPassword && hasSubmit {
			res.HasLogin = true
		}
	}
}
