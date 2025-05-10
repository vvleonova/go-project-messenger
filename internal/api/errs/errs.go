package errs

import (
	"net/http"
)

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func BadRequest(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusUnauthorized, Message: msg}
}

func ServiceUnavailable(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusServiceUnavailable, Message: msg}
}

func NotFound(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusNotFound, Message: msg}
}

func InternalServerError(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusInternalServerError, Message: msg}
}

func Forbidden(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusForbidden, Message: msg}
}
