package models

import (
	"time"
)

type AccessToken struct {
	OwnedEntity
	AccessToken string    `json:"accessToken" gorm:"type:string;"`
	TokenType   string    `json:"tokenType" gorm:"type:string;"`
	ExpiresAt   time.Time `json:"expiresAt" gorm:"type:datetime;"`
}
