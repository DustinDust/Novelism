package main

import (
	"fmt"
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"gin_stuff/internals/models"
	router "gin_stuff/internals/routers"
	"gin_stuff/internals/utils"
	"log"
	"strings"
	"time"

	eLog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"

	"gin_stuff/internals/middlewares"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type Application struct {
	EchoInstance *echo.Echo
	Config       *viper.Viper
	DB           *sqlx.DB
	Models       models.Models
}

func NewApplication() *Application {
	// load config from file
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}

	// create new echo (server) instance
	e := echo.New()

	// create new application
	app := &Application{
		EchoInstance: e,
		Config:       config,
	}

	// apply configuration
	app.EchoInstance.Debug = true
	app.EchoInstance.Logger.SetLevel(eLog.DEBUG)
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
					event = utils.Logger.Error()
					msg = "REQ_ERR"
				} else {
					event = utils.Logger.Info()
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
			event := utils.Logger.Error()
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

	// db configuration
	dbUri := app.Config.GetString("database.uri")
	dbConfig := database.DBConfig{
		MaxIdleConnections: app.Config.GetInt("database.max_db_conns"),
		MaxOpenConnections: app.Config.GetInt("database.max_open_conns"),
		MaxIdleTime:        app.Config.GetDuration("database.max_idle_time"),
	}
	if dbInstance, err := database.OpenDB(dbUri, dbConfig); err != nil {
		app.LogFatalf("Can't open connection to database: %v", err)
	} else {
		log.Printf("Connected to database at %s", strings.Split(dbUri, "@")[1])
		app.DB = dbInstance
	}

	// models
	app.Models = models.NewModels(app.DB)

	// router
	r := router.NewRouter(&app.Models)
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
	jwtRequiredMiddleware := middlewares.NewJwtMiddleware()
	//gloabl prefix
	api := app.EchoInstance.Group("/api")

	// Authentication group
	auth := api.Group("/auth")
	auth.POST("/sign-in", r.Login)
	auth.POST("/sign-up", r.Register)

	//book group
	bookAPI := api.Group("/book", jwtRequiredMiddleware)
	bookAPI.GET("", r.FindBooks)
	bookAPI.GET("/:id", r.GetBook)
	bookAPI.POST("", r.CreateBook)
	bookAPI.PATCH("/:id", r.UpdateBook)
	bookAPI.DELETE("/:id", r.DeleteBook)

	// chapter API
	chapterAPI := api.Group("/book/:bookId/chapter")
	chapterAPI.GET("", r.FindChapters)
	chapterAPI.POST("", r.CreateChapter, jwtRequiredMiddleware)
	chapterAPI.PATCH("/:chapterNo", r.UpdateChapter, jwtRequiredMiddleware)
	chapterAPI.DELETE("/:chapterNo", r.DeleteChapter, jwtRequiredMiddleware)
}
