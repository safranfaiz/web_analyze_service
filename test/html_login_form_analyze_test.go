package test

import (
	"api/analyze"
	"api/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeHtmlLoginForm(t *testing.T) {
	htmlContent := `<!DOCTYPE html> <html lang="en"> <head> <meta charset="UTF-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0"> </head>
	<body> <div class="login-container"> <form action="/login" method="POST"> <div class="form-group"> <label for="username">Username:</label>
	<input type="text" id="username" name="username" required> </div> <div class="form-group"> 
	<label for="password">Password:</label> <input type="password" id="password" name="password" required>
	 </div><button type="submit">Log In</button></form></div></body></html>`

	wc := response.WebContent{
		Content: htmlContent,
	}
	res := response.SuccessResponse{}
	analyze.AnalyzeHtmlLoginForm(&wc, &res)
	assert.Equal(t, true, res.HasLogin, "Web Page Login from analyze testing...")
}
