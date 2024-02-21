package router

import (
	"gin_stuff/internals/models"
	"gin_stuff/internals/services"
	"gin_stuff/internals/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Router struct {
	Model  *models.Models
	Mailer *services.MailerService
}

func NewRouter(model *models.Models) Router {
	return Router{
		Model: model,
	}
}

func (r Router) GetConfig(c echo.Context) error {
	key := c.Param("key")
	userId, err := utils.RetreiveUserIdFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := r.Model.User.Get(int64(userId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(key) <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
		key:    viper.GetViper().Get(key),
	})
}

func (r Router) SendTestMail(c echo.Context) error {
	mailerInfo := struct {
		To    string `json:"to"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}
	if err := c.Bind(&mailerInfo); err != nil {
		utils.Logger.Error().AnErr("binding", err).Send()
		return r.badRequestError(err)
	}
	mail := services.Mail{
		From:    "no-reply@novelism.com",
		To:      mailerInfo.To,
		Subject: mailerInfo.Title,
		Content: mailerInfo.Text,
	}
	if err := r.Mailer.Perform(&mail); err != nil {
		utils.Logger.Error().AnErr("mail_err", err)
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// reutrn http error status 500
func (r Router) serverError(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func (r Router) notFoundError(err error) error {
	return echo.NewHTTPError(http.StatusNotFound, err.Error())
}

func (r Router) badRequestError(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func (r Router) forbiddenError(err error) error {
	return echo.NewHTTPError(http.StatusForbidden, err.Error())
}

func (r Router) unauthorizedError(err error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
}
