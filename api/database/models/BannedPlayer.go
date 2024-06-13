package models

import uuid "github.com/satori/go.uuid"

type BannedPlayer struct {
	GameId   uuid.UUID `json:"-" gorm:"type:uuid;"`
	PlayerId uuid.UUID `json:"-" gorm:"type:uuid;"`
}
