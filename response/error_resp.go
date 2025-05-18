package response

type ErrorResponse struct {
	Message  string `json:"message"`
	ErrorMsg any    `json:"errorMsg"`
	Code     int    `json:"statusCode"`
}

// ErrorResponse function is responsible for create and return a new ErrorResponse.
func ErrorResponseMsg(message string, err any, code int) ErrorResponse {
	return ErrorResponse{
		Message:  message,
		ErrorMsg: err,
		Code:     code,
	}
}
