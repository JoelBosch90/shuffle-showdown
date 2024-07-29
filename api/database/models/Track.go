package models

import (
	"time"
)

type Track struct {
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CheckStartedAt   time.Time  `json:"checkStartedAt"`
	CheckCompletedAt time.Time  `json:"checkCompletedAt"`
	VerifiedAt       time.Time  `json:"verifiedAt"`
	Id               string     `json:"id" gorm:"type:string; primaryKey;"`
	Name             string     `json:"name" gorm:"type:string;"`
	ReleaseYear      uint       `json:"releaseYear" gorm:"type:uint;"`
	Artists          []Artist   `json:"artists" gorm:"many2many:track_artists;"`
	PreviewUrl       string     `json:"previewUrl" gorm:"type:string;"`
	IsPlayable       bool       `json:"isPlayable" gorm:"type:bool;"`
	Playlists        []Playlist `json:"-" gorm:"many2many:playlist_tracks;"`
}
