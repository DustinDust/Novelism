package main

import (
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"gin_stuff/internals/models"
	router "gin_stuff/internals/routers"
	"log"
	"strings"

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
	app.EchoInstance.Logger.SetHeader(viper.GetString("logging.log_header"))
	app.EchoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: viper.GetString("logging.http_log_format"),
	}))
	app.EchoInstance.Use(middleware.Recover())
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
	bookAPI := api.Group("/books", jwtRequiredMiddleware)
	bookAPI.GET("", r.FindBooks)
	bookAPI.GET("/:id", r.GetBook)
	bookAPI.POST("", r.CreateBook)
	bookAPI.PATCH("/:id", r.UpdateBook)
	bookAPI.DELETE("/:id", r.DeleteBook)
}
