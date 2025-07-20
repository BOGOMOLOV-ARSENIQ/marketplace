package database

import (
	"fmt"
	"log"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
)

var DB *gorm.DB

// подключение к бд
func InitDB() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=marketplace_db port=5432 sslmode=disable TimeZone=Europe/Moscow"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully opened!")

	migrate()
}

// автоматическая миграция схем бд
func migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Listing{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed successfully!")
}

// проверка на уникальность
func IsUniqueConstraintError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
