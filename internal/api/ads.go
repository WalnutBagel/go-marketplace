package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
	"github.com/WalnutBagel/go-marketplace/internal/models"
)

func GetAdsHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := middleware.GetUsername(r)
	if !ok || username == "" {
		http.Error(w, "Отсутствует авторизация", http.StatusUnauthorized)
		return
	}

	// Парсим параметры пагинации и сортировки с дефолтами
	page := 1
	limit := 10
	sortField := "created_at"
	order := "DESC"

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		} else {
			http.Error(w, "Невалидный параметр page", http.StatusBadRequest)
			return
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		} else {
			http.Error(w, "Невалидный параметр limit", http.StatusBadRequest)
			return
		}
	}

	if sf := r.URL.Query().Get("sort"); sf != "" {
		allowedSorts := map[string]bool{"created_at": true, "price": true, "title": true}
		if allowedSorts[sf] {
			sortField = sf
		} else {
			http.Error(w, "Невалидное поле сортировки", http.StatusBadRequest)
			return
		}
	}

	if o := strings.ToUpper(r.URL.Query().Get("order")); o == "ASC" || o == "DESC" {
		order = o
	}

	offset := (page - 1) * limit

	var ads []models.Ad
	err := db.GetDB().
		Preload("User").
		Order(sortField + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&ads).Error
	if err != nil {
		http.Error(w, "Ошибка при получении объявлений", http.StatusInternalServerError)
		return
	}

	type AdResp struct {
		ID          uint    `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		Price       float64 `json:"price"`
		CreatedAt   string  `json:"created_at"`
		User        struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}

	resp := make([]AdResp, len(ads))
	for i, ad := range ads {
		resp[i].ID = ad.ID
		resp[i].Title = ad.Title
		resp[i].Description = ad.Description
		resp[i].ImageURL = ad.ImageURL
		resp[i].Price = ad.Price
		resp[i].CreatedAt = ad.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		resp[i].User.ID = ad.User.ID
		resp[i].User.Username = ad.User.Username
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
