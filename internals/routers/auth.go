package router

import (
	"gin_stuff/internals/services"
	"gin_stuff/internals/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r Router) Login(c echo.Context) error {
	payload := LoginPayload{}
	if err := c.Bind(&payload); err != nil {
		return r.badRequestError(err)
	}
	v := utils.NewValidator()
	if err := v.ValidateStruct(&payload); err != nil {
		return r.badRequestError(err)
	}
	user, err := r.Queries.GetUserByUsername(c.Request().Context(), payload.Username)
	if err != nil {
		return r.serverError(err)
	}
	crypt := services.NewCryptoService()
	if err := crypt.Match(payload.Password, user.PasswordHash); err != nil {
		return r.unauthorizedError(err)
	}
	token, err := r.JwtService.SignAccessToken(map[string]interface{}{
		"userId": user.ID,
	})
	if err != nil {
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[LoginData]{OK: true, Data: LoginData{
		AccessToken: token.Token,
		User:        user,
	}})
}
