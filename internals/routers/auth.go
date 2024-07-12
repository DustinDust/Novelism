package router

import (
	"errors"
	"gin_stuff/internals/data"
	"gin_stuff/internals/services"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (r Router) SignIn(c echo.Context) error {
	payload := SignInPayload{}
	if err := c.Bind(&payload); err != nil {
		return r.badRequestError(err)
	}
	if err := r.validator.ValidateStruct(&payload); err != nil {
		return r.badRequestError(err)
	}
	user, err := r.queries.GetUserByUsername(c.Request().Context(), payload.Username)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return r.notFoundError(err)
		default:
			return r.serverError(err)
		}
	}
	crypt := services.NewCryptoService()
	if err := crypt.Match(payload.Password, user.PasswordHash); err != nil {
		return r.unauthorizedError(err)
	}
	token, err := r.jwt.SignAccessToken(map[string]interface{}{
		"userId": user.ID,
	})
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[SignInData]{OK: true, Data: SignInData{
		AccessToken: token.Token,
		User:        user,
	}})
}

func (r Router) SignUp(c echo.Context) error {
	payload := SignUpPayload{}
	if err := c.Bind(&payload); err != nil {
		return r.badRequestError(err)
	}
	if err := r.validator.ValidateStruct(payload); err != nil {
		return r.badRequestError(err)
	}
	crypt := services.NewCryptoService()
	passwordHash, err := crypt.Hash(payload.Password)
	if err != nil {
		return r.serverError(err)
	}
	user, err := r.queries.InsertUser(c.Request().Context(), data.InsertUserParams{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: passwordHash,
		Verified:     pgtype.Bool{Bool: false},
	})
	if err != nil {
		return r.serverError(err)
	}
	token, err := r.jwt.SignAccessToken(map[string]interface{}{
		"userId": user.ID,
	})
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusCreated, Response[SignUpData]{
		OK: true,
		Data: SignUpData{
			AccessToken: token.Token,
			User:        user,
		},
	})
}
