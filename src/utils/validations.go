package utils

import (
	"fmt"
	"regexp"

	"github.com/Numostanley/d8er_app/models"
	"gorm.io/gorm"
)

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
