package config

import "github.com/spf13/viper"

type Config struct {
	PostgresDsn string `mapstructure:"POSTGRES_DSN"`
	JWTSecret	 string `mapstructure:"JWT_SECRET"`
	Environtment string `mapstructure:"ENV"`
	CorsOrigins  string `mapstructure:"CORS_ORIGINS"`
	Domain			 string `mapstructure:"DOMAIN"`
}

func Load() *Config {
	cfg := &Config{
		PostgresDsn: viper.GetString("POSTGRES_DSN"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		Environtment: viper.GetString("ENV"),
		CorsOrigins: viper.GetString("CORS_ORIGINS"),
		Domain: viper.GetString("DOMAIN"),
	}
	return cfg
}
