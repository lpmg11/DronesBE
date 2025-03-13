package config

import "github.com/spf13/viper"

type Config struct {
	PostgresDsn string `mapstructure:"POSTGRES_DSN"`
}

func Load() *Config {
	cfg := &Config{
		PostgresDsn: viper.GetString("POSTGRES_DSN"),
	}
	return cfg
}
