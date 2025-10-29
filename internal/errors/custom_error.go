package errors

import "net/http"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewBadRequest(message string) *APIError {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorized(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewNotFound(message string) *APIError {
	return &APIError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewInternal(message string) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
