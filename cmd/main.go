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
		&models.DroneModel{},
		&models.Drone{},
		&models.Warehouse{},
		&models.Product{},
		&models.Shipment{},
		&models.Client{},
		&models.Budget{},
		&models.BudgetTransaction{},
	)
	if err != nil {
		log.Fatal("Error realizando las migraciones: ", err)
	}

	if err := pg.DB.Exec("CREATE EXTENSION IF NOT EXISTS cube").Error; err != nil {
		log.Fatal("Error creating cube extension: ", err)
	}
	if err := pg.DB.Exec("CREATE EXTENSION IF NOT EXISTS earthdistance").Error; err != nil {
		log.Fatal("Error creating earthdistance extension: ", err)
	}

	log.Println("Migraciones realizadas correctamente")

	log.Println(math.Round(utilities.Distance(14.5995, -90.5153, 14.6225, -90.5135)*100) / 100)

	router := router.Router(cfn, pg)
	router.Run()

}
