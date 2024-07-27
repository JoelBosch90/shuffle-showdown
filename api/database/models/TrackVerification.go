package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TrackVerification struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	TrackId     string    `json:"-" gorm:"type:string;"`
	Track       Playlist  `json:"track" gorm:"foreignKey:TrackId;references:Id;"`
}
