package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Numostanley/d8er_app/db"
	"github.com/Numostanley/d8er_app/models"
	"github.com/Numostanley/d8er_app/serializers"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func ValidatePassword(password string) error {
	passwordRegex := regexp.MustCompile(`[A-Za-z]`)     // Check for at least one letter
	digitRegex := regexp.MustCompile(`\d`)              // Check for at least one digit
	specialCharRegex := regexp.MustCompile(`[@$!%*?&]`) // Check for at least one special character

	if len(password) < 6 || !passwordRegex.MatchString(password) || !digitRegex.MatchString(password) || !specialCharRegex.MatchString(password) {
		return fmt.Errorf("invalid password format")
	}

	return nil
}

func UserExistsByEmail(db *gorm.DB, emailToCheck string) (bool, error) {
	var user models.User
	result := db.Where("email = ?", emailToCheck).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UserExistsByPhone(db *gorm.DB, phoneToCheck string) (bool, error) {
	var user models.User
	result := db.Where("phone_number = ?", phoneToCheck).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UserValidation(user *models.User) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.PhoneNumber == "" {
		return fmt.Errorf("phone_number is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if user.FirstName == "" {
		return fmt.Errorf("first_name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last_name is required")
	}
	return nil
}
