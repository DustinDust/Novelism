package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

func ErrorForbiddenResource(resources ...string) error {
	if len(resources) == 0 {
		return NewError("permissions required to access this resource(s)")
	} else {
		r := strings.Join(resources, ", ")
		return NewError(fmt.Sprintf("permission required: %s", r))
	}
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
