package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WalnutBagel/go-marketplace/internal/api"
	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
	"github.com/WalnutBagel/go-marketplace/internal/models"
)

func main() {
	_, err := db.ConnectWithRetry()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	err = db.GetDB().AutoMigrate(&models.User{}, &models.Ad{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/login", api.LoginHandler)

	// Пример защищённого эндпоинта с выводом имени пользователя из контекста
	protectedHandler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UserContextKey).(string)
		if !ok || username == "" {
			http.Error(w, "Пользователь не найден в контексте", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Доступ разрешён! Привет, %s!", username)
	}))
	mux.Handle("/protected", protectedHandler)

	mux.Handle("/ads", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetAdsHandler(w, r)
		case http.MethodPost:
			api.CreateAdHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
