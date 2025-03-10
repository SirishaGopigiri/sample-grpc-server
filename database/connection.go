package database

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/cloudsqlconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect_to_DB(ctx context.Context, config *Config) (*gorm.DB, error) {
	_, err := cloudsqlconn.NewDialer(ctx)
	if err != nil {
		log.Fatalf("Error creating Cloud SQL dialer: %v", err)
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.DBUser, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	fmt.Println("Successfully connected to Cloud SQL!")

	db.AutoMigrate(&User{})
	return db, nil
}
