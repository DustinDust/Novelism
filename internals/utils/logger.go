package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(zerolog.ConsoleWriter(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: zerolog.TimeFormatUnix,
	FormatTimestamp: func(i interface{}) string {
		rawSec, ok := i.(json.Number)
		if !ok {
			return "INVALID_TIME"
		}
		sec, err := rawSec.Int64()
		if err != nil {
			return "INVALID_TIME"
		}
		time := time.Unix(sec, 0)
		return fmt.Sprintf("%d:%d:%d", time.Local().Hour(), time.Local().Minute(), time.Local().Second())
	},
}))
