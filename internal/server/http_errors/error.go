package http_errors

import "net/http"

type HttpError struct {
	Status  int
	Message string
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewErrBadRequest(m string) *HttpError {
	return &HttpError{
		Status:  http.StatusBadRequest,
		Message: m,
	}
}

func NewErrInternal() *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}
}

func NewErrForbidden(m string) *HttpError {
	return &HttpError{
		Status:  http.StatusForbidden,
		Message: m,
	}
}

func NewErrUnauthorized(m string) *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: m,
	}
}
