package models

import "time"

type Playlist struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Id            string    `json:"id" gorm:"type:string; primaryKey;"`
	Name          string    `json:"name" gorm:"type:string;"`
	CountryCode   string    `json:"-" gorm:"type:string;"`
	LastSongAdded string    `json:"-" gorm:"type:string;"`
	Games         []Game    `json:"-"`
	TracksTotal   uint      `json:"tracksTotal" gorm:"type:uint;"`
	Tracks        []Track   `json:"-" gorm:"many2many:playlist_tracks;"`
}
