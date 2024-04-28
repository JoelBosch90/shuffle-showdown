package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Game struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid; primaryKey;"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Playlist       Playlist  `json:"playlist" gorm:"foreignKey:PlaylistId;"`
	SongsToWin     uint      `json:"songsToWin" gorm:"type:uint;"`
	TitleRequired  bool      `json:"titleRequired" gorm:"type:bool;"`
	ArtistRequired bool      `json:"artistRequired" gorm:"type:bool;"`
	Configured     bool      `json:"configured" gorm:"type:bool;"`
}
