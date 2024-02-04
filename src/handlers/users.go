package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
	"github.com/Numostanley/d8er_app/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err = utils.UserValidation(&newUser)

	if err != nil {
		data.Error = fmt.Sprintf("Validation Error: %s", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	emailExists, _ := utils.UserExistsByEmail(db.Database.DB, newUser.Email)
	if emailExists {
		data.Error = "email already exists"
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	phoneNumberExists, _ := utils.UserExistsByPhone(db.Database.DB, newUser.PhoneNumber)
	if phoneNumberExists {
		data.Error = "phone already exists"
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	err = utils.ValidatePassword(newUser.Password)
	if err != nil {
		data.Error = fmt.Sprintf("Password Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	err = newUser.SetNewPassword(newUser.Password)

	if err != nil {
		data.Error = fmt.Sprintf("Password Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	newUser.CreateUserID()
	user := db.Database.DB.Create(&newUser)

	if user.Error != nil {
		data.Error = fmt.Sprintf("Error creating user: %v", user.Error)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	response := serializers.UserDetailSerializer{}
	data.Data = response.GetUserResponse(&newUser)
	data.Message = "User Created Successfully"

	utils.RespondWithJSON(w, 200, data)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, user models.User) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	response := serializers.UserDetailSerializer{}
	data.Message = "User details retrieved successfully"
	data.Data = response.GetUserResponse(&user)

	utils.RespondWithJSON(w, 200, data)
}
