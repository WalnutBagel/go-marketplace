package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
	"github.com/WalnutBagel/go-marketplace/internal/models"
)

// CreateAdRequest описывает структуру входящих данных для создания или обновления объявления.
type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

// AdResponse описывает структуру JSON-ответа с данными объявления.
type AdResponse struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	ImageURL    string       `json:"image_url"`
	Price       float64      `json:"price"`
	CreatedAt   time.Time    `json:"created_at"`
	User        UserResponse `json:"user"`
}

// UserResponse описывает пользователя в ответе.
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// writeJSON отправляет JSON-ответ с указанным HTTP статусом.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError отправляет JSON-ошибку с сообщением.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// getUsernameFromContext извлекает username из контекста запроса.
func getUsernameFromContext(r *http.Request) (string, error) {
	username, ok := middleware.GetUsername(r)
	if !ok || username == "" {
		return "", errors.New("пользователь не авторизован")
	}
	return username, nil
}

// getUserByUsername загружает пользователя из базы по username.
func getUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("пользователь не найден в БД")
	}
	return &user, nil
}

// validateCreateAdRequest проверяет корректность данных для объявления.
func validateCreateAdRequest(req *CreateAdRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	if len(req.Title) < 3 || len(req.Title) > 100 {
		return errors.New("заголовок должен содержать от 3 до 100 символов")
	}
	if len(req.Description) < 10 || len(req.Description) > 1000 {
		return errors.New("описание должно содержать от 10 до 1000 символов")
	}
	if req.Price <= 0 {
		return errors.New("цена должна быть больше нуля")
	}
	return nil
}

// CreateAdHandler обрабатывает создание нового объявления.
func CreateAdHandler(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "невалидный JSON")
		return
	}

	if err := validateCreateAdRequest(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := getUserByUsername(username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ad := models.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		UserID:      user.ID,
	}

	if err := db.GetDB().Create(&ad).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка при создании объявления")
		return
	}

	resp := AdResponse{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		CreatedAt:   ad.CreatedAt,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	}

	writeJSON(w, http.StatusCreated, resp)
}

// UpdateAdHandler обрабатывает обновление существующего объявления.
func UpdateAdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/ads/")
	adID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный ID объявления")
		return
	}

	username, err := getUsernameFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "невалидный JSON")
		return
	}

	if err := validateCreateAdRequest(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbConn := db.GetDB()
	var ad models.Ad
	if err := dbConn.Preload("User").First(&ad, adID).Error; err != nil {
		writeError(w, http.StatusNotFound, "объявление не найдено")
		return
	}

	if ad.User.Username != username {
		writeError(w, http.StatusForbidden, "нет прав для изменения объявления")
		return
	}

	ad.Title = req.Title
	ad.Description = req.Description
	ad.ImageURL = req.ImageURL
	ad.Price = req.Price

	if err := dbConn.Save(&ad).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка при обновлении объявления")
		return
	}

	resp := AdResponse{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		CreatedAt:   ad.CreatedAt,
		User: UserResponse{
			ID:       ad.User.ID,
			Username: ad.User.Username,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteAdHandler обрабатывает удаление объявления.
func DeleteAdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/ads/")
	adID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неверный ID объявления")
		return
	}

	username, err := getUsernameFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	dbConn := db.GetDB()
	var ad models.Ad
	if err := dbConn.Preload("User").First(&ad, adID).Error; err != nil {
		writeError(w, http.StatusNotFound, "объявление не найдено")
		return
	}

	if ad.User.Username != username {
		writeError(w, http.StatusForbidden, "нет прав для удаления объявления")
		return
	}

	if err := dbConn.Delete(&ad).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка при удалении объявления")
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
