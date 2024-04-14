package models

import (
	"time"
)

type AccessToken struct {
	OwnedEntity
	AccessToken string    `json:"access_token" gorm:"type:string;"`
	TokenType   string    `json:"token_type" gorm:"type:string;"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"type:datetime;"`
}
