package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/models"
	"github.com/WalnutBagel/go-marketplace/internal/services"
	"github.com/WalnutBagel/go-marketplace/internal/utils"
)

// --- STRUCTS ---

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// --- HANDLERS ---

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if err := validateCredentials(req.Username, req.Password); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	var existing models.User
	if err := db.GetDB().Where("username = ?", req.Username).First(&existing).Error; err == nil {
		utils.WriteJSONError(w, http.StatusConflict, "Пользователь с таким логином уже существует")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Ошибка при хэшировании пароля")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Ошибка при сохранении пользователя")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Логин и пароль обязательны")
		return
	}

	var user models.User
	err := db.GetDB().Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	token, err := services.GenerateJWT(user.Username)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// --- HELPERS ---

func validateCredentials(username, password string) error {
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("логин должен быть от 3 до 30 символов")
	}
	if strings.Contains(username, " ") {
		return fmt.Errorf("логин не должен содержать пробелы")
	}
	if len(password) < 6 {
		return fmt.Errorf("пароль должен быть не менее 6 символов")
	}
	return nil
}
