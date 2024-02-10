package routers

import (
	"strings"

	"github.com/Numostanley/auth_app/env"
	"github.com/Numostanley/auth_app/handlers"
	"github.com/Numostanley/auth_app/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func GetRoutes() *chi.Mux {
	enV := env.GetEnv{}
	enV.LoadEnv()
	allowedHosts := enV.AllowedHosts

	mainRouter := chi.NewRouter()
	mainRouter.Use(middlewares.LoggingMiddleware)
	mainRouter.Use(middlewares.AllowedHostsMiddleware(strings.Split(allowedHosts, ",")))

	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	healthRouter := chi.NewRouter()
	healthRouter.Get("/", handlers.HandlerReadiness)

	userRouters := GetUserRouters()
	authRouters := GetAuthRouters()

	v1Router.Mount("/oauth", authRouters)
	v1Router.Mount("/users", userRouters)
	v1Router.Mount("/healthz", healthRouter)
	mainRouter.Mount("/v1", v1Router)

	return mainRouter

}
