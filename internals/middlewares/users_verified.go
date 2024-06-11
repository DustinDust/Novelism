package middlewares

import (
	"fmt"
	"gin_stuff/internals/repositories"
	"gin_stuff/internals/services"
	"gin_stuff/internals/utils"
	"log"

	"github.com/labstack/echo/v4"
)

// need the model object
// pass it versus recreate? Id
func NewUserVerificationRequireMiddleware(repository repositories.IUserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtService := services.JWTService{}
			userId, err := jwtService.RetrieveUserIdFromContext(c)
			if err != nil {
				log.Println(fmt.Errorf("can't find user id inside context object: %v", err))
				return utils.ErrorUnauthorized
			}
			user, err := repository.Get(int64(userId))
			if err != nil {
				log.Println(fmt.Errorf("can't find user: %v", err))
				return utils.ErrorUnauthorized
			}
			if !user.Verified {
				return utils.ErrorUnverfiedUser
			}
			return next(c)
		}
	}
}
