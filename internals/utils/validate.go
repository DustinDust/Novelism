package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	Validator *validator.Validate
}

type StructValidationErrors struct {
	FieldErrors validator.ValidationErrors
}

func NewValidator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Get json name from field using reflect & tags
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Custom validations
	v.RegisterValidation("strongPassword", strongPassword)
	v.RegisterValidation("birthday", birthday)

	return &Validator{
		Validator: v,
	}
}

func birthday(fl validator.FieldLevel) bool {
	birthdayString := fl.Field().String()
	// parse iso timestamp
	birthday, err := time.Parse(time.RFC3339, birthdayString)
	if err != nil {
		return false
	}
	if birthday.UTC().Before(time.Now().UTC()) {
		return true
	}
	return false
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
		verr, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		return &StructValidationErrors{
			FieldErrors: verr,
		}
	}
	return nil
}

// Debug & help subscribe StructValidationError to the Error interface
func (ve StructValidationErrors) Error() string {
	return ve.FieldErrors.Error()
}

// Return a echo.HttpError
func (ve *StructValidationErrors) TranslateToHttpError() error {
	errMessage := "invalid data format"
	errData := []interface{}{}

	for _, e := range ve.FieldErrors {
		errData = append(errData, echo.Map{
			"field": e.Field(),
			"expected": fmt.Sprintf("%s%s", e.ActualTag(), func() string {
				if e.Param() != "" {
					return "=" + e.Param()
				} else {
					return ""
				}
			}()),
			"got":   e.Value(),
			"error": e.Error(),
		})
	}
	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
		"message": errMessage,
		"errors":  errData,
	})
}
