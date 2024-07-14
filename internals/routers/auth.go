package router

import (
	"context"
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
			return r.unauthorizedError(err)
		default:
			return r.serverError(err)
		}
	}
	crypt := services.NewCryptoService()
	if err := crypt.Match(payload.Password, user.PasswordHash); err != nil {
		return r.unauthorizedError(err)
	}
	token, err := r.jwt.SignAccessToken(user.ID)
	if err != nil {
		return r.serverError(err)
	}
	c.SetCookie(&http.Cookie{
		Name:    "novelism_auth",
		Value:   token.Token,
		Path:    "/",
		Expires: token.ExpiresAt,
	})
	return c.JSON(http.StatusOK, Response[SignInData]{OK: true, Data: SignInData{
		User: user,
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
	tx, err := r.db.BeginTx(c.Request().Context(), pgx.TxOptions{})
	if err != nil {
		return r.serverError(err)
	}
	defer tx.Rollback(context.Background())
	user, err := r.queries.WithTx(tx).InsertUser(c.Request().Context(), data.InsertUserParams{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: passwordHash,
		Verified:     pgtype.Bool{Bool: false},
		Status:       data.NullUserStatus{UserStatus: data.UserStatusActive},
	})
	if err != nil {
		return r.badRequestError(err)
	}
	if err := tx.Commit(c.Request().Context()); err != nil {
		return r.serverError(err)
	}
	token, err := r.jwt.SignAccessToken(user.ID)
	if err != nil {
		return r.serverError(err)
	}
	c.SetCookie(&http.Cookie{
		Name:    "novelism_auth",
		Value:   token.Token,
		Path:    "/",
		Expires: token.ExpiresAt,
	})
	return c.JSON(http.StatusCreated, Response[SignUpData]{
		OK: true,
		Data: SignUpData{
			User: user,
		},
	})
}
