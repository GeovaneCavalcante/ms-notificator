package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	envVars *Environments
)

type Environments struct {
	APIPort      string `mapstructure:"API_PORT"`
	Environment  string `mapstructure:"ENVIRONMENT"`
	AwsRegion    string `mapstructure:"AWS_REGION"`
	AwsBaseUrl   string `mapstructure:"AWS_BASE_URL"`
	AwsSnsArn    string `mapstructure:"AWS_SNS_ARN"`
	MongoAddress string `mapstructure:"MONGO_ADDRESS"`
	DbName       string `mapstructure:"DB_NAME"`
}

func LoadEnvVars() *Environments {
	viper.SetConfigFile(".env")
	viper.SetDefault("API_PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "local")
	viper.SetDefault("AWS_REGION", "")
	viper.SetDefault("AWS_BASE_URL", "")
	viper.SetDefault("AWS_SNS_ARN", "")
	viper.SetDefault("MONGO_ADDRESS", "")
	viper.SetDefault("DB_NAME", "")

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
