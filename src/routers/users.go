package routers

import (
	"github.com/Numostanley/auth_app/handlers"
	"github.com/Numostanley/auth_app/middlewares"
	"github.com/go-chi/chi/v5"
)

func GetUserRouters() *chi.Mux {
	userRouter := chi.NewRouter()

	userRouter.Post("/", handlers.CreateUserHandler)
	userRouter.Get("/", middlewares.AuthenticationMiddleware(handlers.GetUserHandler))
	userRouter.Post("/verify_email", handlers.VerifyEmailHandler)
	userRouter.Get("/request_code", handlers.RequestCodeHandler)
	userRouter.Post("/verify_password_change", handlers.VerifyPasswordChangeHandler)
	userRouter.Post("/password_reset", middlewares.AuthenticationMiddleware(handlers.PasswordResetHandler))

	return userRouter
}
