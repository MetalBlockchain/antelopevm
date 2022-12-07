package service

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
