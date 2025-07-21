package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func ConnectWithRetry() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		sslMode := os.Getenv("DB_SSLMODE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)

		for attempts := 1; attempts <= 10; attempts++ {
			dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				sqlDB, err := dbInstance.DB()
				if err == nil && sqlDB.Ping() == nil {
					log.Println("✅ Успешное подключение к БД")
					return
				}
			}

			log.Printf("⌛ Попытка подключения к БД %d/10...", attempts)
			time.Sleep(2 * time.Second)
		}
	})

	if dbInstance == nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	return dbInstance, nil
}

func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("База данных не инициализирована. Вызовите ConnectWithRetry() сначала.")
	}
	return dbInstance
}
