package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresClient struct {
	DB *gorm.DB
}

func NewPostgresClient(dsn string) *PostgresClient {
	log.Println("Connecting to Postgres: " + dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}
	log.Println("Connected to Postgres")

	return &PostgresClient{DB: db}
}
