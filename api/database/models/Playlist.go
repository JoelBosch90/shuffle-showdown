package models

import "time"

type Playlist struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Id            string    `json:"id" gorm:"type:string; primaryKey;"`
	Name          string    `json:"name" gorm:"type:string;"`
	CountryCode   string    `json:"countryCode" gorm:"type:string;"`
	LastSongAdded string    `json:"lastSongAdded" gorm:"type:string;"`
	Games         []Game    `json:"games"`
	TracksTotal   uint      `json:"tracksTotal" gorm:"type:uint;"`
	Tracks        []Track   `json:"tracks" gorm:"many2many:playlist_tracks;"`
}
