package models

import (
	"time"
)

type Track struct {
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Id         string     `json:"id" gorm:"type:string; primaryKey;"`
	Name       string     `json:"name" gorm:"type:string;"`
	Artists    []Artist   `json:"artists" gorm:"many2many:track_artists;"`
	PreviewUrl string     `json:"previewUrl" gorm:"type:string;"`
	IsPlayable bool       `json:"isPlayable" gorm:"type:bool;"`
	Playlists  []Playlist `json:"-" gorm:"many2many:playlist_tracks;"`
	AlbumId    string     `json:"-" gorm:"type:string;"`
	Album      Album      `json:"album" gorm:"foreignKey:AlbumId; references:Id"`
}
