package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Game struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	PlaylistId     string    `json:"-" gorm:"type:string;"`
	Playlist       Playlist  `json:"playlist" gorm:"foreignKey:PlaylistId;references:Id;"`
	IsRunning      bool      `json:"isRunning" gorm:"type:bool;"`
	SongsToWin     uint      `json:"songsToWin" gorm:"type:uint;"`
	TitleRequired  bool      `json:"titleRequired" gorm:"type:bool;"`
	ArtistRequired bool      `json:"artistRequired" gorm:"type:bool;"`
	Configured     bool      `json:"configured" gorm:"type:bool;"`
	OwnerId        uuid.UUID `json:"-" gorm:"type:uuid;"`
	Owner          Player    `json:"owner" gorm:"foreignKey:OwnerId;references:Id;"`
	Players        []Player  `json:"players" gorm:"many2many:game_players;"`
	BannedPlayers  []Player  `json:"-" gorm:"many2many:black_listed_players;"`
}
