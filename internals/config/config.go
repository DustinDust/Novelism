package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	if configDir, found := os.LookupEnv("CONFIG_DIR"); found {
		viper.AddConfigPath(configDir)
	} else {
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./internals/config")
		viper.AddConfigPath(".")
	}

	// set up environment
	env, found := os.LookupEnv("ENV")
	if !found {
		env = "development"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s.toml", env))
	viper.SetConfigType("toml")
	return viper.ReadInConfig()
}
