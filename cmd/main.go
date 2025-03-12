package main

import (
	"drones-be/internal/config"
	"drones-be/internal/models"
	"drones-be/internal/storage"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, proceeding...", err)
	}

	cfn := config.Load()

	pg := storage.NewPostgresClient(cfn.PostgresDsn)

	err := pg.DB.AutoMigrate(
		&models.User{},
		&models.Provider{},
	)
	if err != nil {
		log.Fatal("Error realizando las migraciones: ", err)
	}

}
