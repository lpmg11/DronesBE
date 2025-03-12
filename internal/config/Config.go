package config

import "github.com/spf13/viper"

type Config struct {
	AwsRegion   string `mapstructure:"AWS_REGION"`
	PostgresDsn string `mapstructure:"POSTGRES_DSN"`
}

func Load() *Config {
	cfg := &Config{
		AwsRegion:   viper.GetString("AWS_REGION"),
		PostgresDsn: viper.GetString("POSTGRES_DSN"),
	}
	return cfg
}
