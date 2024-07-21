package app

import (
	"context"
	"fmt"
	"gin_stuff/internals/config"
	router "gin_stuff/internals/routers"
	"gin_stuff/internals/services"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Application struct {
	engine *echo.Echo
	db     *pgx.Conn
}

func NewApplication() *Application {
	// load config from file
	config.LoadConfig()
	loggerService := services.NewLoggerService(os.Stdout)

	// create new echo (server) instance
	e := echo.New()

	// create new application
	app := &Application{
		engine: e,
	}

	// apply configuration
	app.engine.Debug = true
	app.engine.Use(middleware.RequestID())
	app.engine.Use(middleware.RequestLoggerWithConfig(
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
	app.engine.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
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
	app.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.engine.Static("", "assets")

	// db configuration
	dbUri := viper.GetString("database.uri")
	if dbInstance, err := OpenDB(); err != nil {
		log.Fatalf("Can't open connection to database: %v", err)
	} else {
		log.Printf("Connected to database at %s", strings.Split(dbUri, "@")[1])
		app.db = dbInstance
	}

	// handle shut down of stuff
	app.engine.Server.RegisterOnShutdown(func() {
		err := app.db.Close(context.Background())
		if err != nil {
			loggerService.LogError(err, "fail to gracefully shutdown database connection")
		}
	})

	r, err := router.New(app.db)
	if err != nil {
		loggerService.LogFatal(err, "failed to init router")
	}
	app.RegisterRoute(r)

	return app
}

// Register the routes in server
func (app Application) RegisterRoute(r *router.Router) {
	api := app.engine.Group("/api")
	auth := api.Group("/auth")
	book := api.Group("/book")

	auth.POST("/sign-in", r.SignIn)
	auth.POST("/sign-up", r.SignUp)

	book.GET("", r.BrowseBooks)
	book.POST("", r.CreateBook, r.JWTMiddleware("access"))
	book.PATCH("/:bookId", r.UpdateBook, r.JWTMiddleware("access"))
	book.DELETE("/:bookId", r.DeleteBook, r.JWTMiddleware("access"))
}

func (app Application) Run() error {
	addr := viper.GetString("general.server")
	if addr == "" {
		addr = ":80"
	}
	return app.engine.Start(addr)
}

func OpenDB() (*pgx.Conn, error) {
	uri := viper.GetString("database.uri")
	if pgxConn, err := pgx.Connect(context.Background(), uri); err != nil {
		return nil, err
	} else {
		return pgxConn, nil
	}
}
