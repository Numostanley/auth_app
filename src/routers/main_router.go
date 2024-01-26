package routers

import (
	"os"
	"strings"

	"github.com/Numostanley/d8er_app/handlers"
	"github.com/Numostanley/d8er_app/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func GetRoutes() *chi.Mux {
	godotenv.Load()
	allowedHosts := os.Getenv("ALLOWED_HOSTS")

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
