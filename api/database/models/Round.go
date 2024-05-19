package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Round struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	Number    uint      `json:"number" gorm:"type:uint;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	GameId    uuid.UUID `json:"-" gorm:"type:uuid;"`
	Game      Game      `json:"-" gorm:"foreignKey:GameId;references:Id;"`
	TrackId   string    `json:"-" gorm:"type:string;"`
	Track     Track     `json:"track" gorm:"foreignKey:TrackId;references:Id;"`
	PlayerId  uuid.UUID `json:"playerId" gorm:"type:uuid;"`
	Player    Player    `json:"-" gorm:"foreignKey:PlayerId;references:Id;"`
}
