package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccessToken struct {
	Id          uuid.UUID `json:"id" gorm:"column:id; type:uuid; primary_key;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AccessToken string    `json:"accessToken" gorm:"column:access_token; type:string;"`
	TokenType   string    `json:"tokenType" gorm:"column:token_type; type:string;"`
	ExpiresAt   time.Time `json:"expiresAt" gorm:"column:expires_at; type:datetime;"`
}
