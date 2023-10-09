package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("strongPassword", strongPassword)

	return &Validator{
		Validator: v,
	}
}

func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check if the password has at least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check if the password contains at least one uppercase letter, one lowercase letter, and one digit
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	for _, char := range password {
		if 'A' <= char && char <= 'Z' {
			hasUppercase = true
		}
		if 'a' <= char && char <= 'z' {
			hasLowercase = true
		}
		if '0' <= char && char <= '9' {
			hasDigit = true
		}
	}

	return hasUppercase && hasLowercase && hasDigit
}

func (v *Validator) ValidateStruct(args interface{}) error {
	if err := v.Validator.Struct(args); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
