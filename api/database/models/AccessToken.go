package models

import (
	"time"
)

type AccessToken struct {
	OwnedEntity
	AccessToken string    `json:"accessToken" gorm:"column:access_token; type:string;"`
	TokenType   string    `json:"tokenType" gorm:"column:token_type; type:string;"`
	ExpiresAt   time.Time `json:"expiresAt" gorm:"column:expires_at; type:datetime;"`
}
