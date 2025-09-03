package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

func NewGormDB(cfg *config.Config) *gorm.DB {
	log.Printf("Database config: StatsStorage=%s, URL=%s", cfg.Database.StatsStorage, cfg.Database.URL)
	
	if cfg.Database.StatsStorage == "memory" {
		log.Println("Using memory storage, skipping database connection")
		return nil
	}

	log.Println("Connecting to PostgreSQL...")
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	log.Println("Running database migration...")
	err = db.AutoMigrate(&model.StatsEntry{})
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}

	log.Println("Database connection and migration successful")
	return db
}