package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Numostanley/auth_app/internal/serializers"
	"github.com/Numostanley/auth_app/internal/utils"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	data := serializers.ResponseSerializer{
		Success: true,
		Message: "",
		Data:    struct{}{},
		Error:   "",
	}

	var requestMap map[string]string
	oauthConfig := utils.OauthConfig{}
	oauthConfig.Initialize()

	authentication := utils.TokenAuthentication{}

	err := json.NewDecoder(r.Body).Decode(&requestMap)
	if err != nil {
		data.Error = fmt.Sprintf("Error parsing JSON: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	clientID := requestMap["client_id"]
	clientSecret := requestMap["client_secret"]
	grantType := requestMap["grant_type"]
	email := requestMap["email"]
	password := requestMap["password"]
	scope := requestMap["scope"]

	_, user, err := utils.PerformAuthentication(clientID, clientSecret, grantType, email, password, scope)
	if err != nil {
		data.Error = fmt.Sprintf("Authentication Error: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	token, err := authentication.GetOauthToken(scope, *user)

	if err != nil {
		data.Error = fmt.Sprintf("Error generting token: %v", err)
		data.Success = false
		utils.RespondWithError(w, 400, data)
		return
	}

	responseData := serializers.TokenResponseSerializer{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   oauthConfig.TokenExpiryTime,
		Scope:       scope,
	}

	data.Data = responseData
	data.Message = "User login Successful"

	utils.RespondWithJSON(w, 200, data)
}
