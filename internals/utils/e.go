package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func NewError(message string, code int) error {
	return echo.NewHTTPError(code, message)
}
