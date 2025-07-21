package main

import (
	"log"
	"net/http"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/models"
	"github.com/WalnutBagel/go-marketplace/internal/router"
)

func main() {
	_, err := db.ConnectWithRetry()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("✅ Успешное подключение к БД")

	if err := db.GetDB().AutoMigrate(&models.User{}, &models.Ad{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", router.NewRouter()))
}
