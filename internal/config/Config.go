package config

import "github.com/spf13/viper"

type Config struct {
	PostgresDsn string `mapstructure:"POSTGRES_DSN"`
	JWTSecret	 string `mapstructure:"JWT_SECRET"`
}

func Load() *Config {
	cfg := &Config{
		PostgresDsn: viper.GetString("POSTGRES_DSN"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}
	return cfg
}
