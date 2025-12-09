package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/krishna102001/dependecy-injection/internal/handlers"
)

func SetupRoutes(hd *handlers.Handler, logger *slog.Logger, r *chi.Mux) http.Handler {
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// r.Get("/",hd)

	return r
}
