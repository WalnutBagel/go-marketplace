package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/services"
)

type contextKey string

const UserContextKey = contextKey("username")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует авторизация", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат авторизации", http.StatusUnauthorized)
			return
		}

		username, err := services.ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Кладем username в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
