package utils

import "errors"

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorUnauthorized       = errors.New("unauthorized")
	ErrorRecordsNotFound    = errors.New("record(s) not found")
	ErrorValidationStruct   = errors.New("invalid data format")
	ErrorInvalidRouteParam  = errors.New("invalid route parameters")
	ErrorForbiddenResource  = errors.New("permissions required to access this resource(s)")
	ErrorInvalidQueryParams = errors.New("invalid route query")
	ErrorInvalidModel       = errors.New("invalid model object")
)
