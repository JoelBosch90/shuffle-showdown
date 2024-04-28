package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccessToken struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AccessToken string    `json:"accessToken" gorm:"type:string;"`
	TokenType   string    `json:"tokenType" gorm:"type:string;"`
	ExpiresAt   time.Time `json:"expiresAt" gorm:"type:datetime;"`
}
