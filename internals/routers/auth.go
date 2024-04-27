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

type ForgetPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordPayload struct {
	UserId      int    `json:"userId" validate:"required,min=1"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"password" validate:"required,min=6,max=20,strongPassword"`
}

type LoginResponseData struct {
	AccessToken  services.SignedJwtResult `json:"accessToken"`
	RefreshToken services.SignedJwtResult `json:"refreshToken,omitempty"`
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
	accessToken, err := r.JwtService.SignAccessToken(user.ID)
	if err != nil {
		return r.serverError(err)
	}
	refreshToken, err := r.JwtService.SignRefreshToken(user.ID)
	if err != nil {
		return r.serverError(err)
	}
	user.RefreshToken = refreshToken.Token
	user.RefreshTokenValidUntil = &refreshToken.ExpiresAt
	r.Model.User.Update(user)
	return c.JSON(http.StatusOK, Response[LoginResponseData]{
		OK: true,
		Data: LoginResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
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
		return c.JSON(http.StatusOK, Response[any]{
			OK: true,
		})
	}
	return r.unauthorizedError(utils.ErrorInvalidToken)
}

func (r Router) ResendVerificationEmail(c echo.Context) error {
	userId, err := r.JwtService.RetreiveUserIdFromContext(c)
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
	cryptoService := services.NewCryptoService()
	verificationToken := cryptoService.GenerateSecureToken(32)
	user.VerificationToken = verificationToken
	if err := r.Model.User.Update(user); err != nil {
		return r.serverError(err)
	}
	if err := r.MailerService.Perform(&services.Mail{
		From:    "no-reply@novelism.com",
		To:      user.Email,
		Subject: "Please verify your email!",
		Content: fmt.Sprintf("https://frontent-link/verify-email?token=%s&user_id=%d", user.VerificationToken, user.ID),
	}); err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[any]{
		OK: true,
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
	cryptoService := services.NewCryptoService()
	verificationToken := cryptoService.GenerateSecureToken(32)
	user.VerificationToken = verificationToken
	if err := r.Model.User.Insert(user); err != nil {
		return r.badRequestError(err)
	}
	if err := r.MailerService.Perform(&services.Mail{
		From:    "no-reply@novelism.com",
		To:      user.Email,
		Subject: "Welcome to novelism! Please verify your email",
		Content: fmt.Sprintf("https://frontend-link/verify-email?token=%s&user_id=%d", user.VerificationToken, user.ID),
	}); err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusCreated, Response[any]{
		OK: true,
	})
}

func (r Router) ForgetPassword(c echo.Context) error {
	validate := utils.NewValidator()
	payload := new(ForgetPasswordPayload)
	if err := c.Bind(payload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(payload); err != nil {
		return r.badRequestError(err)
	}
	user, err := r.Model.User.GetByEmail(payload.Email, "active")
	if err != nil {
		return r.badRequestError(err)
	}
	cryptoService := services.NewCryptoService()
	passwordResetToken := cryptoService.GenerateSecureToken(32)
	user.PasswordResetToken = passwordResetToken
	if err := r.Model.User.Update(user); err != nil {
		return r.serverError(err)
	}

	if err := r.MailerService.Perform(&services.Mail{
		From:    "no-reply@novelism.com",
		To:      user.Email,
		Subject: "Please reset your password",
		Content: fmt.Sprintf("https://frontend-link/reset-password?token=%s&user_id=%d", user.PasswordResetToken, user.ID),
	}); err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[any]{
		OK: true,
	})
}

func (r Router) ResetPassword(c echo.Context) error {
	validate := utils.NewValidator()
	payload := new(ResetPasswordPayload)
	if err := c.Bind(payload); err != nil {
		return r.badRequestError(err)
	}
	if err := validate.ValidateStruct(payload); err != nil {
		return r.badRequestError(err)
	}
	user, err := r.Model.User.Get(int64(payload.UserId))
	if err != nil {
		return r.badRequestError(err)
	}
	if payload.Token != user.PasswordResetToken {
		return r.unauthorizedError(utils.NewError("invalid error", 403, nil))
	}
	if err := user.SetPassword(payload.NewPassword); err != nil {

		return r.serverError(err)
	}
	if err := r.Model.User.Update(user); err != nil {
		return r.serverError(err)
	}
	return c.JSON(200, Response[any]{
		OK: true,
	})
}

func (r Router) Me(c echo.Context) error {
	userId, err := r.JwtService.RetreiveUserIdFromContext(c)
	if err != nil {
		return r.unauthorizedError(err)
	}
	user, err := r.Model.User.Get(int64(userId))
	if err != nil {
		return r.badRequestError(err)
	}
	return c.JSON(200, Response[models.User]{
		OK:   true,
		Data: *user,
	})
}
