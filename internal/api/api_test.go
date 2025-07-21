package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WalnutBagel/go-marketplace/internal/db"
	"github.com/WalnutBagel/go-marketplace/internal/models"
	"github.com/WalnutBagel/go-marketplace/internal/router"
)

func TestMain(m *testing.M) {
	// Установка переменных окружения, например JWT_SECRET
	os.Setenv("JWT_SECRET", "testsecret")

	// Подключение к базе данных перед тестами
	_, err := db.ConnectWithRetry()
	if err != nil {
		panic("Ошибка подключения к БД в тестах: " + err.Error())
	}

	code := m.Run()

	// Очистка таблиц после всех тестов
	db.GetDB().Exec("DELETE FROM users")
	db.GetDB().Exec("DELETE FROM ads")

	os.Exit(code)
}

func TestRegisterHandler(t *testing.T) {
	router := router.NewRouter()

	// Подготовка тела запроса
	payload := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	// Создаем HTTP-запрос
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ в recorder
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидали статус 200 OK, получили: %d", w.Code)
	}

	var respUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &respUser); err != nil {
		t.Fatalf("Ошибка разбора JSON ответа: %v", err)
	}

	if respUser.Username != "testuser" {
		t.Errorf("Ожидали username 'testuser', получили: %s", respUser.Username)
	}
}
