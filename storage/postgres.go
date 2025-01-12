package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBUser     string
	DBHost     string
	DBName     string
	DBPassword string
	DBPort     string
	DBSSLMode  string
	DebugMode  bool
}

func NewPostgreStorage(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, "UTC")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if cfg.DebugMode {
		db = db.Debug()
	}
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	return db, nil
}
