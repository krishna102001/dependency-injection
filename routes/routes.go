package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/krishna102001/dependecy-injection/internal/handlers"
	middlewares "github.com/krishna102001/dependecy-injection/internal/middleware"
)

func SetupRoutes(handler *handlers.Handler, logger *slog.Logger, r *chi.Mux) http.Handler {
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.DefaultLogger)

	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.IsAuthenticated(logger))
		r.Get("/", handler.GetPlainText)
	})

	return r
}
