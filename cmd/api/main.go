package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	app := NewApplication()

	go func() {
		if err := app.EchoInstance.Start(app.Config.GetString("general.server")); err != nil && err != http.ErrServerClosed {
			app.LogFatalf("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.EchoInstance.Shutdown(ctx); err != nil {
		app.LogFatalf("%v", err)
	}
}
