package utils

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(zerolog.ConsoleWriter(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: zerolog.TimeFormatUnix,
}))
