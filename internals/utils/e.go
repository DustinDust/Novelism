package utils

import "net/http"

type NovelismError struct {
	Message string
	Code    int
	Data    any
}

var (
	ErrorInvalidCredentials = NewError("invalid credentials", http.StatusUnauthorized)
	ErrorUnauthorized       = NewError("unauthorized", http.StatusUnauthorized)
	ErrorRecordsNotFound    = NewError("record(s) not found", http.StatusNotFound)
	ErrorValidationStruct   = NewError("invalid data format", http.StatusBadRequest)
	ErrorInvalidRouteParam  = NewError("invalid route parameters", http.StatusBadRequest)
	ErrorForbiddenResource  = NewError("permissions required to access this resource(s)", http.StatusForbidden)
	ErrorInvalidQueryParams = NewError("invalid route query", http.StatusBadRequest)
	ErrorInvalidModel       = NewError("invalid model object", http.StatusBadRequest)
	ErrorUnverfiedUser      = NewError("unverified user", http.StatusUnauthorized)
	ErrorInvalidToken       = NewError("invalid token", http.StatusUnauthorized)
)

func NewError(message string, code int, data ...any) NovelismError {
	return NovelismError{
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func (e NovelismError) Error() string {
	return e.Message
}
