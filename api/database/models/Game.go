package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Game struct {
	Id             uuid.UUID    `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	PlaylistId     string       `json:"-" gorm:"type:string;"`
	Playlist       Playlist     `json:"playlist" gorm:"foreignKey:PlaylistId;references:Id;"`
	HasStarted     bool         `json:"hasStarted" gorm:"type:bool;"`
	HasFinished    bool         `json:"hasFinished" gorm:"type:bool;"`
	SongsToWin     uint         `json:"songsToWin" gorm:"type:uint;"`
	TitleRequired  bool         `json:"titleRequired" gorm:"type:bool;"`
	ArtistRequired bool         `json:"artistRequired" gorm:"type:bool;"`
	Configured     bool         `json:"configured" gorm:"type:bool;"`
	OwnerId        uuid.UUID    `json:"-" gorm:"type:uuid;"`
	Owner          Player       `json:"owner" gorm:"foreignKey:OwnerId;references:Id;"`
	Players        []Player     `json:"players" gorm:"many2many:game_players;"`
	BannedPlayers  []Player     `json:"-" gorm:"many2many:black_listed_players;"`
	Rounds         []Round      `json:"rounds" gorm:"foreignKey:GameId;references:Id;"`
	GamePlayers    []GamePlayer `json:"-" gorm:"foreignKey:GameId;references:Id;"`
	WonTracks      []WonTrack   `json:"wonTracks" gorm:"foreignKey:GameId;references:Id;"`
}
