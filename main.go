package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/routers"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}

	db.InitDB()
	// utils.SeedClient()

	mainRouter := routers.GetRoutes()

	server := &http.Server{
		Handler: mainRouter,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	log.Fatal(server.ListenAndServe())
}
