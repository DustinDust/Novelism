package router

import (
	"fmt"
	"gin_stuff/internals/repositories"
	"gin_stuff/internals/services"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Router struct {
	Repository    *repositories.Repository
	MailerService *services.MailerService
	JwtService    *services.JWTService
	LoggerService *services.LoggerService
	GenAIService  *services.GeminiService
}

func NewRouter(repository *repositories.Repository, mailerService *services.MailerService, loggerService *services.LoggerService, genAIService *services.GeminiService) Router {
	return Router{
		Repository:    repository,
		MailerService: mailerService,
		LoggerService: loggerService,
		GenAIService:  genAIService,
		JwtService:    &services.JWTService{}, // recreate each router creation since it does not initiate any object instance
	}
}

type Response[T interface{}] struct {
	OK       bool                  `json:"ok"`
	Data     T                     `json:"data"`
	Metadata repositories.Metadata `json:"metadata"`
}

// route handler to test runtime config
func (r Router) GetConfig(c echo.Context) error {
	key := c.Param("key")
	userId, err := r.JwtService.RetrieveUserIdFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := r.Repository.User.Get(int64(userId))
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

// route handler to test email functionality
func (r Router) SendTestMail(c echo.Context) error {
	mailerInfo := struct {
		To    string `json:"to"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}
	if err := c.Bind(&mailerInfo); err != nil {
		log.Printf("Error binding body %v", err)
		return r.badRequestError(err)
	}
	mail := services.Mail{
		From:    "no-reply@novelism.com",
		To:      mailerInfo.To,
		Subject: mailerInfo.Title,
		Content: mailerInfo.Text,
	}
	if err := r.MailerService.Perform(&mail); err != nil {
		log.Println(err)
		return r.serverError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

func (r Router) TestFileUpload(c echo.Context) error {
	name := c.FormValue("name")
	file, err := c.FormFile("image")
	if err != nil {
		return r.serverError(err)
	}
	src, err := file.Open()
	if err != nil {
		return r.serverError(fmt.Errorf("first check: %+v", err))
	}
	dst, err := os.Create(fmt.Sprintf("%s-%s", name, file.Filename))
	if err != nil {
		return r.serverError(fmt.Errorf("second check: %+v", err))
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return r.serverError(fmt.Errorf("third check: %+v", err))
	}
	return c.JSON(200, Response[string]{
		OK:   true,
		Data: file.Filename,
	})
}

func (r Router) TestAIPrompt(c echo.Context) error {
	type PromptBody struct {
		Prompt string `json:"prompt"`
	}
	body := new(PromptBody)
	if err := c.Bind(body); err != nil {
		return r.badRequestError(err)
	}
	data, err := r.GenAIService.GenerateText(body.Prompt)
	if err != nil {
		r.serverError(err)
	}
	return c.JSON(http.StatusOK, Response[[]string]{
		OK:   true,
		Data: data,
	})
}

// return http errors
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
