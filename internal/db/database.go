package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // глобальная переменная для хранения подключения

// ConnectWithRetry подключается к БД с повторными попытками
func ConnectWithRetry() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)

	var err error
	for attempts := 1; attempts <= 10; attempts++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := DB.DB()
			if err != nil {
				return nil, err
			}
			err = sqlDB.Ping()
			if err == nil {
				log.Println("✅ Успешное подключение к БД")
				return DB, nil
			}
		}

		log.Printf("Пытаемся подключиться к БД... попытка %d/10", attempts)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
}

// GetDB возвращает текущее подключение к базе
func GetDB() *gorm.DB {
	return DB
}
