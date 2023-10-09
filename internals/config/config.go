package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func LoadConfig() (*viper.Viper, error) {
	// set up environment
	env, found := os.LookupEnv("ENV")
	if !found {
		env = "development"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s.toml", env))
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return viper.GetViper(), nil
}
