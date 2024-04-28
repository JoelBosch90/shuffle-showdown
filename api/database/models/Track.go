package models

import (
	"time"
)

type Track struct {
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Id           string    `json:"id" gorm:"type:string; primaryKey;"`
	Name         string    `json:"name" gorm:"type:string;"`
	ReleaseYear  uint      `json:"releaseYear" gorm:"type:uint;"`
	ReleaseMonth uint      `json:"releaseMonth" gorm:"type:uint;"`
	ReleaseDay   uint      `json:"releaseDay" gorm:"type:uint;"`
	Artists      []Artist  `json:"artists" gorm:"many2many:track_artists;"`
}
