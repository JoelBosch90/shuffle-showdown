package models

import uuid "github.com/satori/go.uuid"

type GamePlayer struct {
	GameId   uuid.UUID `json:"-" gorm:"type:uuid;"`
	PlayerId uuid.UUID `json:"-" gorm:"type:uuid;"`
	Order    uint      `json:"order" gorm:"type:uint;"`
}
