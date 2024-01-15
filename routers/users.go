package routers

import (
	"github.com/Numostanley/d8er_app/handlers"
	"github.com/Numostanley/d8er_app/middlewares"
	"github.com/go-chi/chi"
)

func GetUserRouters() *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Post("/", handlers.CreateUserHandler)
	userRouter.Get("/", middlewares.AuthenticationMiddleware(handlers.GetUserHandler))

	return userRouter
}