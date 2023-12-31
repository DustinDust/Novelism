package router

import (
	"errors"
	"gin_stuff/internals/models"
	"gin_stuff/internals/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginPayload struct {
	Username          string `json:"username" validate:"required,min=6"`
	PlaintextPassword string `json:"password" validate:"required,min=6,max=20"`
}

type RegisterPayload struct {
	Username          string `json:"username" validate:"required,min=6"`
	PlaintextPassword string `json:"password" validate:"required,min=6,max=20,strongPassword"`
	Email             string `json:"email" validate:"required,email"`
}

// Handler
func (r Router) Login(c echo.Context) error {
	validate := utils.NewValidator()
	loginPayload := new(LoginPayload)
	if err := c.Bind(loginPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.ValidateStruct(loginPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	user, err := r.Model.User.Login(loginPayload.Username, loginPayload.PlaintextPassword)
	if err != nil {
		if errors.Is(err, utils.ErrorInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	accessToken, err := utils.SignToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"accessToken": accessToken})
}

func (r Router) Register(c echo.Context) error {
	validate := utils.NewValidator()
	registerPayload := new(RegisterPayload)
	if err := c.Bind(registerPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.ValidateStruct(registerPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	user := &models.User{
		Username: registerPayload.Username,
		Email:    registerPayload.Email,
		Status:   "active",
	}
	if err := user.SetPassword(registerPayload.PlaintextPassword); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if err := r.Model.User.Insert(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	accessToken, err := utils.SignToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"accessToken": accessToken})
}
