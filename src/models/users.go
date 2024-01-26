package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *User) CreateUserID() {
	u.ID = uuid.New()
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (u *User) ValidatePassword(password string) bool {
	err := ComparePasswords(u.Password, password)
	return err == nil
}

func (u *User) HasActiveSession() bool {
	lastLogin := u.LastLogin
	if lastLogin != nil {
		timeDifference := time.Since(*lastLogin)
		return timeDifference < 30*time.Minute
	}
	return false
}

func (u *User) SetNewPassword(password string) error {
	password, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = password
	return nil
}

func (u *User) ValidateUserAgainstClientID(clientID string) bool {
	if u.IsAdmin {
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
