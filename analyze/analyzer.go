package analyze

import "api/response"

// Analyzer defines the interface for all web content analyzers.
type Analyzer interface {
	Analyze(wc *response.WebContent, res *response.SuccessResponse) *response.ErrorResponse
}
