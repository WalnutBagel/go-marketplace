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
	username, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || username == "" {
		http.Error(w, "Отсутствует авторизация", http.StatusUnauthorized)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortField := r.URL.Query().Get("sort")
	order := strings.ToUpper(r.URL.Query().Get("order"))

	if pageStr == "" {
		pageStr = "1"
	}
	if limitStr == "" {
		limitStr = "10"
	}
	if sortField == "" {
		sortField = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "Невалидный параметр page", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		http.Error(w, "Невалидный параметр limit", http.StatusBadRequest)
		return
	}

	allowedSortFields := map[string]bool{"created_at": true, "price": true, "title": true}
	if !allowedSortFields[sortField] {
		http.Error(w, "Невалидное поле сортировки", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit

	var ads []models.Ad
	err = db.GetDB().
		Order(sortField + " " + order).
		Limit(limit).
		Offset(offset).
		Preload("User").
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
