package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
	"github.com/WalnutBagel/go-marketplace/internal/models"
)

type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

func CreateAdHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || username == "" {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	var req CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	// Обрезаем пробелы
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	// Валидация
	if len(req.Title) < 3 || len(req.Title) > 100 ||
		len(req.Description) < 10 || len(req.Description) > 1000 ||
		req.Price <= 0 {
		http.Error(w, "Невалидные данные", http.StatusBadRequest)
		return
	}

	// Ищем пользователя
	var user models.User
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "Пользователь не найден в БД", http.StatusInternalServerError)
		return
	}

	// Создаем объявление
	ad := models.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		UserID:      user.ID,
	}

	if err := db.GetDB().Create(&ad).Error; err != nil {
		http.Error(w, "Ошибка при создании объявления", http.StatusInternalServerError)
		return
	}

	// Ответ
	resp := map[string]interface{}{
		"id":          ad.ID,
		"title":       ad.Title,
		"description": ad.Description,
		"image_url":   ad.ImageURL,
		"price":       ad.Price,
		"created_at":  ad.CreatedAt,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
