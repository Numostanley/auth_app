package models

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type VerificationCode struct {
	*gorm.Model
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Code    string    `json:"code" gorm:"type:text;not null"`
	IsValid bool      `json:"is_valid" gorm:"default:True"`
}

func (vCode *VerificationCode) VerifyCodeValidity() bool {
	date := &vCode.UpdatedAt
	if date != nil {
		timeDifference := time.Since(*date)
		return timeDifference < 30*time.Minute
	}
	return false
}

func GetCode(db *gorm.DB, code string) (*VerificationCode, error) {
	var vCode VerificationCode
	result := db.Where("code = ?", code).First(&vCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &vCode, fmt.Errorf("invalid code")
	}
	return &vCode, nil
}

func CreateVerificationCode(user *User, db *gorm.DB) *VerificationCode {
	vCode := VerificationCode{UserID: user.ID}
	result := db.Where("user_id = ?", user.ID).First(&vCode)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		vCode.Code = GenerateOTP()
		vCode.IsValid = true
		db.Create(&vCode)
		return &vCode
	} else {
		vCode.Code = GenerateOTP()
		vCode.IsValid = true
		db.Save(&vCode)
		return &vCode
	}
}

type User struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	FirstName       string     `json:"first_name" gorm:"type:text;not null"`
	LastName        string     `json:"last_name" gorm:"type:text;not null"`
	Email           string     `json:"email" gorm:"type:text;not null;uniqueIndex"`
	Password        string     `json:"password" gorm:"type:text;not null"`
	LastLogin       *time.Time `json:"last_login"`
	Gender          *string    `json:"gender"`
	IsAdmin         bool       `json:"is_admin" gorm:"default:false"`
	IsEmailVerified bool       `json:"is_email_verified" gorm:"default:false"`
	PhoneNumber     string     `json:"phone_number" gorm:"type:text;not null;uniqueIndex"`
}

func (user *User) CreateUserID() {
	user.ID = uuid.New()
}

func (user *User) CreateUser(db *gorm.DB) error {
	err := UserValidation(user)
	if err != nil {
		return err
	}

	emailExists, _ := UserExistsByEmail(db, user.Email)
	if emailExists {
		err = fmt.Errorf("email already exists")
		return err
	}

	phoneNumberExists, _ := UserExistsByPhone(db, user.PhoneNumber)
	if phoneNumberExists {
		err = fmt.Errorf("phone already exists")
		return err
	}

	err = ValidatePassword(user.Password)
	if err != nil {
		err = fmt.Errorf("password Error: %v", err)
		return err
	}

	err = user.SetNewPassword(user.Password)
	if err != nil {
		err = fmt.Errorf("password Error: %v", err)
		return err
	}

	user.CreateUserID()
	newUser := db.Create(&user)
	if newUser.Error != nil {
		err = fmt.Errorf("error creating user: %v", newUser.Error)
		return err
	}
	return nil
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) ValidatePassword(password string) bool {
	err := ComparePasswords(user.Password, password)
	return err == nil
}

func (user *User) SetNewPassword(password string) error {
	password, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = password
	return nil
}

func (user *User) ValidateUserAgainstClientID(clientID string) bool {
	if user.IsAdmin {
		for _, client := range []string{AppClients.AdminAppClient, AppClients.MobileAppClient} {
			if clientID == client {
				return true
			}
		}
		return false
	} else {
		return clientID == AppClients.MobileAppClient
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err
}

func UserExistsByEmail(db *gorm.DB, emailToCheck string) (bool, error) {
	var user User
	result := db.Where("email = ?", emailToCheck).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UserExistsByPhone(db *gorm.DB, phoneToCheck string) (bool, error) {
	var user User
	result := db.Where("phone_number = ?", phoneToCheck).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UserValidation(user *User) error {
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

func GetUserByID(userID uuid.UUID, db *gorm.DB) (*User, error) {
	user := User{ID: userID}
	fetchedUser := db.Where("id = ?", userID).First(&user)

	if fetchedUser.Error != nil {
		return nil, fmt.Errorf("error returning user %s", fetchedUser.Error)
	}
	return &user, nil
}

func GetUserByEmail(email string, db *gorm.DB) (*User, error) {
	user := User{Email: email}
	fetchedUser := db.Where("email = ?", email).First(&user)

	if fetchedUser.Error != nil {
		return nil, fmt.Errorf("error returning user %s", fetchedUser.Error)
	}
	return &user, nil
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

func GenerateOTP() string {
	rand.NewSource(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)
}
