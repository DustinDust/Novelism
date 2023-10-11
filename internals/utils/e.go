package utils

import "errors"

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorUnauthorized       = errors.New("unauthorized")
	ErrorRecordsNotFound    = errors.New("record(s) not found")
	ErrorValidationStruct   = errors.New("invalid data format")
)
