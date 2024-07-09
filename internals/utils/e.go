package utils

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stderr)

func NewError(message string) error {
	err := errors.New(message)
	// log wrapper
	logger.Error().Err(err).Timestamp().Send()
	return err
}

func ErrInvalidCredentials() error {
	return NewError("invalid credentials")
}

func ErrUnauthorized() error {
	return NewError("unauthorized")
}

func ErrorRecordsNotFound() error {
	return NewError("record(s) not found")
}

func ErrorValidationStruct() error {
	return NewError("invalid data format")
}

func ErrorInvalidRouteParam() error {
	return NewError("invalid route parameters")
}

func ErrorForbiddenResource() error {
	return NewError("permissions required to access this resource(s)")
}

func ErrorInvalidQueryParams() error {
	return NewError("invalid route query")
}

func ErrorInvalidModel() error {
	return NewError("invalid model object")
}

func ErrorUnverfiedUser() error {
	return NewError("unverified user")
}

func ErrorInvalidToken() error {
	return NewError("invalid token")
}

func ErrInvalidContext() error {
	return NewError("invalid server context")
}
