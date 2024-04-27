package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Game struct {
	Id             uuid.UUID `json:"id" gorm:"column:id; type:uuid; primary_key;"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	PlaylistId     string    `json:"playlistId" gorm:"column:playlist_id; type:string;"`
	SongsToWin     uint      `json:"songsToWin" gorm:"column:songs_to_win; type:uint;"`
	TitleRequired  bool      `json:"titleRequired" gorm:"column:title_required; type:bool;"`
	ArtistRequired bool      `json:"artistRequired" gorm:"column:artist_required; type:bool;"`
	Configured     bool      `json:"configured" gorm:"column:configured; type:bool;"`
}
