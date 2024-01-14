package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
)

func CreateClient() {
	filename := "extras/clients.json"

	file, err := OpenFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer CloseFile(file)

	clientParams := []models.Client{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&clientParams)
	if err != nil {
		fmt.Println("error decoding json: ", err)
	}

	db := &db.Database.DB

	for _, client := range clientParams {
		models.CreateClient(*db, &client)
		fmt.Println(
			client,
		)
	}
}
