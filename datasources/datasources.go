package datasources

import (
	"fmt"
	"log"
	"os"
	
	"darkoo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Ds struct {
	DB *gorm.DB
}



func InitDS() (*Ds, error) {
	log.Print("Initializing data sources\n")

	DATABASE_HOST := os.Getenv("DATABASE_HOST")
	DATABASE_USER := os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	DATABASE_DB := os.Getenv("DATABASE_DB")
	DATABASE_PORT := os.Getenv("DATABASE_PORT")
	DB_SSL_MODE := os.Getenv("DB_SSL_MODE")

	log.Printf("Connecting to pstgres sql\n")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD, DATABASE_DB, DATABASE_PORT, DB_SSL_MODE)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Print("Error opening database")
		return nil, fmt.Errorf("Error opening database %w", err)
	}

	if err := db.AutoMigrate(
		&models.Group{}, &models.Message{}, &models.User{}, &models.UserGroup{},
	); err != nil {
		log.Print("Error migrating models")
		return nil, fmt.Errorf("Error migrating models: %w", err)
	}
	return &Ds{
		DB: db,
	}, nil
}