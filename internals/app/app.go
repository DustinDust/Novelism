package app

import (
	"fmt"
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"gin_stuff/internals/middlewares"
	"gin_stuff/internals/repositories"
	router "gin_stuff/internals/routers"
	"gin_stuff/internals/services"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Application struct {
	EchoInstance *echo.Echo
	DB           *sqlx.DB
}

func NewApplication() *Application {
	// load config from file
	config.LoadConfig()
	loggerService := services.NewLoggerService()

	// create new echo (server) instance
	e := echo.New()

	// create new application
	app := &Application{
		EchoInstance: e,
	}

	// apply configuration
	app.EchoInstance.Debug = true
	app.EchoInstance.Use(middleware.RequestID())
	app.EchoInstance.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogURI:          true,
			LogStatus:       true,
			LogError:        true,
			LogLatency:      true,
			LogProtocol:     true,
			LogMethod:       true,
			LogRequestID:    true,
			LogResponseSize: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				var event *zerolog.Event
				var msg string
				if v.Error != nil {
					event = loggerService.Logger.Error()
					msg = "REQ_ERR"
				} else {
					event = loggerService.Logger.Info()
					msg = "REQ_OK"
				}
				event.Time("time", v.StartTime.Local())
				event.Str("req_id", v.RequestID)
				event.Str("method", v.Method)
				event.Str("uri", v.URI)
				event.Int("status", v.Status)
				event.Dur("latency", v.Latency)
				event.Str("prot", v.Protocol)
				event.Int64("resp_size", v.ResponseSize)
				event.Stack().Err(v.Error)
				event.Msg(msg)
				return nil
			},
		},
	))
	app.EchoInstance.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			event := loggerService.Logger.Error()
			event.Time("time", time.Now())
			event.Stack().Err(err)
			event.Str("uri", c.Path())
			event.Send()
			fmt.Println(string(stack))
			return err
		},
	}))
	app.EchoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.EchoInstance.Static("/", "assets")

	// db configuration
	dbUri := viper.GetString("database.uri")
	dbConfig := database.DBConfig{
		MaxIdleConnections: viper.GetInt("database.max_db_conns"),
		MaxOpenConnections: viper.GetInt("database.max_open_conns"),
		MaxIdleTime:        viper.GetDuration("database.max_idle_time"),
	}
	if dbInstance, err := database.OpenDB(dbUri, dbConfig); err != nil {
		app.LogFatalf("Can't open connection to database: %v", err)
	} else {
		log.Printf("Connected to database at %s", strings.Split(dbUri, "@")[1])
		app.DB = dbInstance
	}

	// mailer
	mailer, err := services.NewMailerService(services.MailerSMTPConfig{
		Host:     viper.GetString("mailer.host"),
		Port:     viper.GetInt64("mailer.port"),
		Login:    viper.GetString("mailer.login"),
		Password: viper.GetString("mailer.password"),
		Timeout:  viper.GetDuration("mailer.timeout"),
	})
	if err != nil {
		loggerService.LogFatal(err, "fail to initialize mailer service")
	}

	// handle shut down of stuff
	app.EchoInstance.Server.RegisterOnShutdown(func() {
		err = app.DB.Close()
		if err != nil {
			loggerService.LogError(err, "fail to gracefully shutdown database connection")
		}
	})

	repo := repositories.New(app.DB)
	r := router.New(&repo, mailer, &loggerService)
	app.RegisterRoute(r)

	return app
}

// some helper functions for application

// Log fatal error
func (app *Application) LogFatalf(format string, args ...interface{}) {
	app.EchoInstance.Logger.Fatalf(format, args)
}

// Log server information
func (app *Application) LogInfof(format string, args ...interface{}) {
	app.EchoInstance.Logger.Infof(format, args)
}

// Register the routes in server
func (app Application) RegisterRoute(r router.Router) {
	requireAccessToken := middlewares.NewJWTMiddleware("access")
	requireUserVerification := middlewares.NewUserVerificationRequireMiddleware(r.Repository.User)
	//gloabl prefix
	api := app.EchoInstance.Group("/api")

	// Authentication group
	auth := api.Group("/auth")
	auth.POST("/sign-in", r.Login)
	auth.POST("/sign-up", r.Register)
	auth.POST("/verify-email", r.VerifyEmail)
	auth.POST("/resend-verification-mail", r.ResendVerificationEmail, requireAccessToken)
	auth.POST("/forget-password", r.ForgetPassword)
	auth.POST("/reset-password", r.ResetPassword)
	auth.GET("/me", r.Me, requireAccessToken)

	//book group
	bookAPI := api.Group("/book")
	bookAPI.GET("", r.FindBooks, requireAccessToken)
	bookAPI.GET("/:id", r.GetBook)
	bookAPI.POST("", r.CreateBook, requireAccessToken, requireUserVerification)
	bookAPI.PATCH("/:id", r.UpdateBook, requireAccessToken, requireUserVerification)
	bookAPI.DELETE("/:id", r.DeleteBook, requireAccessToken, requireUserVerification)

	// chapter API
	chapterAPI := bookAPI.Group("/:bookId/chapter")
	chapterAPI.GET("", r.FindChapters)
	chapterAPI.POST("", r.CreateChapter, requireAccessToken, requireUserVerification)
	chapterAPI.PATCH("/:chapterNo", r.UpdateChapter, requireAccessToken, requireUserVerification)
	chapterAPI.DELETE("/:chapterNo", r.DeleteChapter, requireAccessToken, requireUserVerification)

	// chapter content
	contentAPI := api.Group("/chapter/:chapterUID/content")
	contentAPI.GET("", r.GetContent, requireAccessToken)
}

func (app Application) Run(addr string) error {
	return app.EchoInstance.Start(addr)
}
