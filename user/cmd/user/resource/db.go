package resource

import (
	"fmt"
	"log"
	"user/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func IntDb(cfg *config.Config) *gorm.DB {
	// Debug print to verify loaded config values
	dsn := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failted to connect to DB:: %v", err)
	}
	log.Println("Connected to DB successfully")
	return db

}
