package serializers

import (
	"time"

	"github.com/Numostanley/auth_app/models"
	"github.com/google/uuid"
)

type UserDetailSerializer struct {
	ID              uuid.UUID  `json:"id"`
	FullName        string     `json:"full_name"`
	Gender          *string    `json:"gender"`
	Email           string     `json:"email"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLogin       *time.Time `json:"last_login"`
	PhoneNumber     string     `json:"phone_number"`
	IsEmailVerified bool       `json:"is_email_verified"`
}

func (userRes *UserDetailSerializer) GetUserResponse(user *models.User) UserDetailSerializer {
	response := UserDetailSerializer{
		ID:              user.ID,
		FullName:        user.GetFullName(),
		Gender:          user.Gender,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		LastLogin:       user.LastLogin,
		PhoneNumber:     user.PhoneNumber,
		IsEmailVerified: user.IsEmailVerified,
	}

	return response
}
