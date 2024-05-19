package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Player struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" gorm:"type:varchar(255);"`
	// Secret is a UUID that is used to identify the player in the game.
	// It should not be exposed to the client, hence it's excluded from the JSON.
	Secret uuid.UUID `json:"-" gorm:"type:uuid;"`
	// No need to expose information about linked games either.
	ParticipatedGames []Game     `json:"-" gorm:"many2many:game_players;"`
	WonTracks         []WonTrack `json:"awardedTracks" gorm:"foreignKey:PlayerId;references:Id;"`
}
