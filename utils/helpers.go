package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

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

func GetOauthToken(scope string, user models.User) (string, error) {

	oauthConfig := OauthConfig{}
	oauthConfig.Initialize()

	key, err := oauthConfig.GetPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error reading private key: %v", err)
	}
	aud := oauthConfig.Audience
	now := time.Now()

	tokenExpiryTime, err := oauthConfig.GetTokenExpiryTime()
	if err != nil {
		return "", fmt.Errorf("error reading tokenExpiryTime key: %v", err)
	}

	payload := map[string]interface{}{
		"iss":          oauthConfig.Issuer,
		"sub":          user.ID,
		"full_name":    user.GetFullName(),
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"scope":        scope,
		"iat":          now,
		"aud":          aud,
		"exp":          now.Add(time.Second * time.Duration(tokenExpiryTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %v", err)
	}

	user.LastLogin = &now
	db := db.Database.DB
	db.Save(&user)

	return tokenString, nil
}

// func ComparePasswords(hashedPassword, inputPassword string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
// 	return err
// }
