package router

import (
	"errors"
	"fmt"
	"gin_stuff/internals/models"
	"gin_stuff/internals/services"
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
		return r.badRequestError(err)
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
			return r.unauthorizedError(err)
		} else {
			return r.serverError(err)
		}
	}
	accessToken, err := utils.JWT.SignToken(user.ID, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"accessToken": accessToken})
}

func (r Router) VerifyEmail(c echo.Context) error {
	validate := utils.NewValidator()
	payload := new(struct {
		Token  string `json:"token" validate:"required"`
		UserID int    `json:"userId" validate:"required"`
	})
	if err := c.Bind(payload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(payload); err != nil {
		return r.badRequestError(err)
	}
	user, err := r.Model.User.Get(int64(payload.UserID))
	if err != nil {
		return r.badRequestError(err)
	}
	if user.Verified || user.Status != "active" {
		return r.badRequestError(fmt.Errorf("invalid user status"))
	}
	if user.VerificationToken == payload.Token {
		user.VerificationToken = ""
		user.Verified = true
		err := r.Model.User.Update(user)
		if err != nil {
			return r.serverError(err)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"ok": true,
		})
	}
	return r.unauthorizedError(utils.ErrorInvalidToken)
}

func (r Router) ResendVerificationEmail(c echo.Context) error {
	userId, err := utils.JWT.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.unauthorizedError(err)
	}
	user, err := r.Model.User.Get(int64(userId))
	if err != nil {
		return r.unauthorizedError(err)
	}
	if user.Verified {
		return r.badRequestError(fmt.Errorf("user is already verified"))
	}
	verificationToken := utils.Crypto.GenerateSecureToken(32)
	user.VerificationToken = verificationToken
	if err := r.Model.User.Update(user); err != nil {
		return r.serverError(err)
	}
	if err := r.Mailer.Perform(&services.Mail{
		From:    "no-reply@novelism.com",
		To:      user.Email,
		Subject: "Please verify your email!",
		Content: fmt.Sprintf("https://frontent-link/verify-email?token=%s&user_id=%d", user.VerificationToken, user.ID),
	}); err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

func (r Router) Register(c echo.Context) error {
	validate := utils.NewValidator()
	registerPayload := new(RegisterPayload)
	if err := c.Bind(registerPayload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(registerPayload); err != nil {
		if verr, ok := err.(*utils.StructValidationErrors); ok {
			return verr.TranslateError()
		} else {
			return r.serverError(err)
		}
	}
	user := &models.User{
		Username: registerPayload.Username,
		Email:    registerPayload.Email,
		Status:   "active",
		Verified: false,
	}
	if err := user.SetPassword(registerPayload.PlaintextPassword); err != nil {
		return r.serverError(err)
	}

	verificationToken := utils.Crypto.GenerateSecureToken(32)
	user.VerificationToken = verificationToken
	if err := r.Model.User.Insert(user); err != nil {
		return r.badRequestError(err)
	}
	if err := r.Mailer.Perform(&services.Mail{
		From:    "no-reply@novelism.com",
		To:      user.Email,
		Subject: "Welcome to novelism! Please verify your email",
		Content: fmt.Sprintf("https://frontent-link/verify-email?token=%s&user_id=%d", user.VerificationToken, user.ID),
	}); err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusCreated, echo.Map{"ok": true})
}
