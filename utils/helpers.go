package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
	"github.com/google/uuid"
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
	w.Write(data)
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
		file.Close()
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

	clientParams := []models.Client{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&clientParams)
	if err != nil {
		log.Println("error decoding json: ", err)
	}

	db := &db.Database.DB
	for _, client := range clientParams {
		_, err := GetClientByClientID(client.ClientID)

		if err != nil {
			models.CreateClient(*db, &client)
			log.Println(
				client,
			)
		}
	}
}

func GetUserByID(userID uuid.UUID) (*models.User, error) {
	user := models.User{ID: userID}
	fetchedUser := db.Database.DB.Where("id = ?", userID).First(&user)

	if fetchedUser.Error != nil {
		return nil, fmt.Errorf("error returning user %s", fetchedUser.Error)
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := models.User{Email: email}
	fetchedUser := db.Database.DB.Where("email = ?", email).First(&user)

	if fetchedUser.Error != nil {
		return nil, fmt.Errorf("error returning user %s", fetchedUser.Error)
	}
	return &user, nil
}

func GetClientByClientID(clientID string) (*models.Client, error) {
	client := models.Client{ClientID: clientID}
	fetchedClient := db.Database.DB.Where("client_id = ?", clientID).First(&client)

	if fetchedClient.Error != nil {
		return nil, fmt.Errorf("error returning client %s", fetchedClient.Error)
	}
	return &client, nil
}
