package service

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
