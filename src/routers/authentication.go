package routers

import (
	"github.com/Numostanley/d8er_app/handlers"
	"github.com/go-chi/chi"
)

func GetAuthRouters() *chi.Mux {
	authRouter := chi.NewRouter()
	authRouter.Post("/token", handlers.AuthHandler)

	return authRouter
}
