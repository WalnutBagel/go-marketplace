package router

import (
	"net/http"
	"strings"

	"github.com/WalnutBagel/go-marketplace/internal/api"
	"github.com/WalnutBagel/go-marketplace/internal/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/login", api.LoginHandler)
	mux.Handle("/ads", middleware.AuthMiddleware(http.HandlerFunc(adRouter)))
	mux.Handle("/ads/", middleware.AuthMiddleware(http.HandlerFunc(adRouter)))

	return mux
}

func adRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/ads" {
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
