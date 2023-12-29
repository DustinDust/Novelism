package utils

import (
	"github.com/rs/zerolog/log"
)

func LogWarning(any map[string]interface{}) {
	logEvent := log.Warn()
	for key, value := range any {
		logEvent.Any(key, value)
	}
	logEvent.Send()
}

func LogInfo(any map[string]interface{}) {
	logEvent := log.Info()
	for key, value := range any {
		logEvent.Any(key, value)
	}
	logEvent.Send()

}
