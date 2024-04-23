package utils

import "net/http"

type NovelismError struct {
	Message string
	Code    int
	Data    any
}

var (
	ErrorInvalidCredentials = NewError("invalid credentials", http.StatusUnauthorized, nil)
	ErrorUnauthorized       = NewError("unauthorized", http.StatusUnauthorized, nil)
	ErrorRecordsNotFound    = NewError("record(s) not found", http.StatusNotFound, nil)
	ErrorValidationStruct   = NewError("invalid data format", http.StatusBadRequest, nil)
	ErrorInvalidRouteParam  = NewError("invalid route parameters", http.StatusBadRequest, nil)
	ErrorForbiddenResource  = NewError("permissions required to access this resource(s)", http.StatusForbidden, nil)
	ErrorInvalidQueryParams = NewError("invalid route query", http.StatusBadRequest, nil)
	ErrorInvalidModel       = NewError("invalid model object", http.StatusBadRequest, nil)
	ErrorUnverfiedUser      = NewError("unverified user", http.StatusUnauthorized, nil)
	ErrorInvalidToken       = NewError("invalid token", http.StatusUnauthorized, nil)
)

func NewError(message string, code int, data any) NovelismError {
	return NovelismError{
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func (e NovelismError) Error() string {
	return e.Message
}
