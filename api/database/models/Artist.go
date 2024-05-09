package models

import "time"

type Artist struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Id        string    `json:"id" gorm:"type:string; primaryKey;"`
	Name      string    `json:"name" gorm:"type:string;"`
	Tracks    []Track   `json:"tracks" gorm:"many2many:track_artists;"`
}
