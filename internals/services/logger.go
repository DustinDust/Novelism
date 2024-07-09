package services

import (
	"io"

	"github.com/rs/zerolog"
)

type LoggerService struct {
	Logger zerolog.Logger
}

func NewLoggerService(out io.Writer) LoggerService {
	return LoggerService{
		Logger: zerolog.New(out),
	}
}

func (service LoggerService) LogError(err error, message string) {
	service.Logger.Error().Err(err).Timestamp().Msg(message)
}

func (service LoggerService) LogInfo(data any, message string) {
	service.Logger.Info().Timestamp().Any("log_data", data).Msg(message)
}

func (service LoggerService) LogDebug(data any, message string) {
	service.Logger.Debug().Timestamp().Any("log_data", data).Msg(message)
}

func (service LoggerService) LogFatal(err error, message string) {
	service.Logger.Fatal().Timestamp().Err(err).Msg(message)
}
