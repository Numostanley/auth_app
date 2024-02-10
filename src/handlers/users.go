package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Numostanley/auth_app/db"
	"github.com/Numostanley/auth_app/models"
	"github.com/Numostanley/auth_app/serializers"
	"github.com/Numostanley/auth_app/utils"
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

func GetUserHandler(w http.ResponseWriter, _ *http.Request, user models.User) {
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

func RequestCodeHandler(w http.ResponseWriter, r *http.Request) {
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

	database := db.Database.DB
	email := r.URL.Query().Get("email")
	if email == "" {
		data.Error = "email is required"
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
	data.Message = "Verification code sent to user"

	go utils.VerificationEmail(user)
	utils.RespondWithJSON(w, 200, data)
}

func VerifyPasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
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
	password := requestMap["password"]

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
	err = user.SetNewPassword(password)
	if err != nil {
		data.Error = fmt.Sprintf("Error setting new password: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	database.Save(&user)
	database.Model(&vCode).Update("is_valid", false)

	data.Message = "User Password changed successfully"
	utils.RespondWithJSON(w, 200, data)
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request, user models.User) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	var requestMap map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestMap)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	database := db.Database.DB
	password := requestMap["password"]
	password2 := requestMap["password2"]

	if password != password2 {
		data.Error = "Passwords must match!!!"
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	err = models.ValidatePassword(password)
	if err != nil {
		data.Error = fmt.Sprintf("Error validating password: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	err = user.SetNewPassword(password)
	if err != nil {
		data.Error = fmt.Sprintf("Error setting new password: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}
	database.Save(&user)

	data.Message = "User Password reset successfully"
	utils.RespondWithJSON(w, 200, data)
}
