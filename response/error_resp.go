package response

type ErrorResponse struct {
	Message  string `json:"message"`
	ErrorMsg string `json:"errorMsg"`
}

// ErrorResponse function is responsible for create and return a new ErrorResponse.
func ErrorResponseMsg(message string, err error) ErrorResponse {
	return ErrorResponse{
		Message:  message,
		ErrorMsg: err.Error(),
	}
}
