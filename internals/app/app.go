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
	EchoInstance *echo.Echo
	DB           *pgx.Conn
}

func NewApplication() *Application {
	// load config from file
	config.LoadConfig()
	loggerService := services.NewLoggerService(os.Stdout)

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
	if dbInstance, err := OpenDB(); err != nil {
		log.Fatalf("Can't open connection to database: %v", err)
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
		log.Fatal(err, "fail to initialize mailer service")
	}

	// handle shut down of stuff
	app.EchoInstance.Server.RegisterOnShutdown(func() {
		err = app.DB.Close(context.Background())
		if err != nil {
			loggerService.LogError(err, "fail to gracefully shutdown database connection")
		}
	})

	r := router.New(mailer, &loggerService)
	app.RegisterRoute(r)

	return app
}

// Register the routes in server
func (app Application) RegisterRoute(r router.Router) {
}

func (app Application) Run() error {
	addr := viper.GetString("general.server")
	if addr == "" {
		addr = ":80"
	}
	return app.EchoInstance.Start(addr)
}

func OpenDB() (*pgx.Conn, error) {
	uri := viper.GetString("database.uri")
	if pgxConn, err := pgx.Connect(context.Background(), uri); err != nil {
		return nil, err
	} else {
		return pgxConn, nil
	}
}
