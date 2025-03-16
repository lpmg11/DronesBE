package main

import (
	"drones-be/internal/config"
	"drones-be/internal/models"
	router "drones-be/internal/routes"
	"drones-be/internal/storage"
	"drones-be/internal/utilities"
	"log"
	"math"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, proceeding...", err)
	}

	cfn := config.Load()

	pg := storage.NewPostgresClient(cfn.PostgresDsn)

	err := pg.DB.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.Drone{},
	)
	if err != nil {
		log.Fatal("Error realizando las migraciones: ", err)
	}

	log.Println("Migraciones realizadas correctamente")

	log.Println(math.Round(utilities.Distance(14.5995, -90.5153, 14.6225, -90.5135)*100) / 100)

	router := router.Router(cfn, pg)
	router.Run()

}
