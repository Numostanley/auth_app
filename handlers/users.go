package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
	"github.com/Numostanley/d8er_app/utils"
	"github.com/google/uuid"
)

func HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	var newUser models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	password, ok := utils.HashPassword(newUser.Password)

	if ok != nil {
		data.Error = fmt.Sprintf("Password Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	newUser.Password = password
	user := db.Database.DB.Create(&newUser)

	if user.Error != nil {
		data.Error = fmt.Sprintf("Error creating user: %v", user.Error)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	response := serializers.UserDetailSerializer{}
	data.Data = response.GetUserResponse(newUser)
	data.Message = "User Created Successfully"

	utils.RespondWithJSON(w, 200, data)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	id := r.URL.Query().Get("id")
	userID, err := uuid.Parse(id)

	if err != nil {
		data.Error = fmt.Sprintf("Error parsing uuid: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	user := models.User{ID: userID}
	fetchedUser := db.Database.DB.Where("id = ?", userID).First(&user)

	if fetchedUser.Error != nil {
		data.Error = fmt.Sprintf("Error getting user: %v", fetchedUser.Error)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	response := serializers.UserDetailSerializer{}
	data.Message = "User details retrieved successfully"
	data.Data = response.GetUserResponse(user)

	utils.RespondWithJSON(w, 200, data)
}
