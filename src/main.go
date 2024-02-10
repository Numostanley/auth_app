package main

import (
	"log"
	"net/http"

	"github.com/Numostanley/auth_app/db"
	"github.com/Numostanley/auth_app/env"
	"github.com/Numostanley/auth_app/routers"
	"github.com/Numostanley/auth_app/utils"
)

func main() {
	enV := env.GetEnv{}
	enV.LoadEnv()

	db.InitDB()
	utils.SeedClient()

	mainRouter := routers.GetRoutes()

	server := &http.Server{
		Handler: mainRouter,
		Addr:    ":" + enV.PortString,
	}

	log.Printf("Server starting on port %v", enV.PortString)
	log.Fatal(server.ListenAndServe())
}
