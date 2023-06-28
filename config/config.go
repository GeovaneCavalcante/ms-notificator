package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	envVars *Environments
)

type Environments struct {
	APIPort string `mapstructure:"API_PORT"`
}

func LoadEnvVars() *Environments {
	viper.SetConfigFile(".env")
	viper.SetDefault("API_PORT", "8080")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Print("unable find or read configuration file: %w", err)
	}

	if err := viper.Unmarshal(&envVars); err != nil {
		fmt.Print("unable to unmarshal configurations from environment: %w", err)
	}

	return envVars
}

func GetEnvVars() *Environments {
	return envVars
}
