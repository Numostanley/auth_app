package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	FirstName    string     `json:"first_name" gorm:"text;not null"`
	LastName     string     `json:"last_name" gorm:"text;not null"`
	Email        string     `json:"email" gorm:"not null;uniqueIndex"`
	Password     string     `json:"password" gorm:"not null"`
	LastLogin    *time.Time `json:"last_login"`
	Gender       *string    `json:"gender"`
	IsSuperAdmin bool       `json:"is_super_admin" gorm:"default:false"`
	PhoneNumber  string     `json:"phone_number" gorm:"text"`
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
