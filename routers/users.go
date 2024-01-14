package routers

import (
	"github.com/Numostanley/d8er_app/handlers"
	"github.com/go-chi/chi"
)

func GetUserRouters() *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Post("/", handlers.HandlerCreateUser)
	userRouter.Get("/", handlers.GetUserHandler)

	return userRouter
}
