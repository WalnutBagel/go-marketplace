package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/services"
)

type contextKey string

const userKey contextKey = "username"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
			return
		}

		username, err := services.ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "Недействительный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Получить имя пользователя из контекста
func GetUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(userKey).(string)
	return username, ok
}
