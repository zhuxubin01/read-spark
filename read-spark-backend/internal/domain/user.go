package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Phone     string    `gorm:"uniqueIndex;size:20" json:"phone"`
	Email     *string   `gorm:"uniqueIndex;size:255" json:"email,omitempty"`
	Nickname  *string   `gorm:"size:100" json:"nickname,omitempty"`
	AvatarURL *string   `gorm:"size:500" json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type UserLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
