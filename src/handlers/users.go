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

	userClientID := models.AppClients.MobileAppClient
	err := utils.BasicAuthentication(r, userClientID)
	if err != nil {
		data.Error = fmt.Sprintf("Authorization Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	var newUser models.User

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUser)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	database := db.Database.DB
	err = newUser.CreateUser(database)

	if err != nil {
		data.Error = fmt.Sprintf("Error creating user: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	response := serializers.UserDetailSerializer{}
	data.Data = response.GetUserResponse(&newUser)
	data.Message = "User Created Successfully"

	go utils.VerificationEmail(&newUser)

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

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	userClientID := models.AppClients.MobileAppClient
	err := utils.BasicAuthentication(r, userClientID)
	if err != nil {
		data.Error = fmt.Sprintf("Authorization Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	var requestMap map[string]string
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&requestMap)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	database := db.Database.DB
	code := requestMap["code"]
	email := requestMap["email"]

	vCode, err := models.GetCode(database, code)
	if err != nil {
		data.Error = fmt.Sprintf("Error validating code: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	if !vCode.VerifyCodeValidity() || !vCode.IsValid {
		data.Error = "Code already used or expired"
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	user, err := models.GetUserByEmail(email, database)
	if err != nil {
		data.Error = fmt.Sprintf("Error validating user: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	user.IsEmailVerified = true
	database.Save(&user)

	database.Model(&vCode).Update("is_valid", false)

	response := serializers.UserDetailSerializer{}
	data.Message = "Email verified successfully"
	data.Data = response.GetUserResponse(user)

	utils.RespondWithJSON(w, 200, data)
}
