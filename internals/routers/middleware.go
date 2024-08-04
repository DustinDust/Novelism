package router

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// @param{strategy}: "access"
func (r Router) JWTMiddleware(strategy string) echo.MiddlewareFunc {
	jwtSecret := viper.GetString(fmt.Sprintf("jwt.%s_secret", strategy))

	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  "user",
		TokenLookup: "header:Authorization:Bearer ,cookie:novelism_auth",
		SigningKey:  []byte(jwtSecret),
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			customClaims := new(struct {
				Claims int `json:"claims"`
				jwt.RegisteredClaims
			})
			token, err := jwt.ParseWithClaims(auth, customClaims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, r.forbiddenError(errors.New("unexpected signing method"))
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				return nil, err
			}
			content, ok := token.Claims.(*struct {
				Claims int `json:"claims"`
				jwt.RegisteredClaims
			})
			if !ok {
				return nil, r.forbiddenError(errors.New("unexpected jwt format"))
			}
			// content.Claims is userId
			user, err := r.queries.GetUserByID(c.Request().Context(), int32(content.Claims))
			if err != nil {
				return nil, r.forbiddenError(err)
			}
			return user, nil
		},
	})
}
