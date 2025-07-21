package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/api"
	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
	"github.com/WalnutBagel/go-marketplace/internal/models"
)

func main() {
	// Подключение к базе данных с автоматической миграцией
	_, err := db.ConnectWithRetry()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("✅ Успешное подключение к БД")

	if err := db.GetDB().AutoMigrate(&models.User{}, &models.Ad{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/login", api.LoginHandler)

	// Единый обработчик для всех /ads и /ads/:id
	mux.Handle("/ads", middleware.AuthMiddleware(http.HandlerFunc(adRouter)))
	mux.Handle("/ads/", middleware.AuthMiddleware(http.HandlerFunc(adRouter)))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func adRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if r.URL.Path == "/ads" {
		switch r.Method {
		case http.MethodGet:
			api.GetAdsHandler(w, r)
		case http.MethodPost:
			api.CreateAdHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
		return
	}

	// Если путь начинается с "/ads/", значит работа с конкретным объявлением
	if strings.HasPrefix(path, "/ads/") {
		switch r.Method {
		case http.MethodPut:
			api.UpdateAdHandler(w, r)
		case http.MethodDelete:
			api.DeleteAdHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}
