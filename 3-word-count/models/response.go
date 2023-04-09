package models

type ErrorResponse struct {
	ErrorTitle   string `json:"error_title"`
	ErrorMessage string `json:"error_message"`
}
