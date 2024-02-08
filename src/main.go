package main

import (
	"log"
	"net/http"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/env"
	"github.com/Numostanley/d8er_app/routers"
	"github.com/Numostanley/d8er_app/utils"
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
