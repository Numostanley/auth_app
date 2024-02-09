package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
)

func RespondWithError(w http.ResponseWriter, code int, data serializers.ResponseSerializer) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", data.Error)
	}
	RespondWithJSON(w, code, data)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func OpenFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CloseFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			return
		}
	}
}

func SeedClient() {
	filename := "extras/clients.json"

	file, err := OpenFile(filename)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer CloseFile(file)

	var clientParams []models.Client

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&clientParams)
	if err != nil {
		log.Println("error decoding json: ", err)
	}

	database := db.Database.DB
	for _, client := range clientParams {
		_, err := models.GetClientByClientID(client.ClientID, database)

		if err != nil {
			err := models.CreateClient(database, &client)
			if err != nil {
				return
			}
			log.Println(
				client,
			)
		}
	}
}
