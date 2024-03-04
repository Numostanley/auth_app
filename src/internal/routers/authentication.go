package routers

import (
	"github.com/Numostanley/auth_app/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetAuthRouters() *chi.Mux {
	authRouter := chi.NewRouter()
	authRouter.Post("/token", handlers.AuthHandler)

	return authRouter
}
