package models

import "time"

type Track struct {
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Id           string    `json:"id" gorm:"column:id; type:string; primary_key;"`
	Name         string    `json:"name" gorm:"column:name; type:string;"`
	ReleaseYear  uint      `json:"releaseYear" gorm:"column:release_year; type:uint;"`
	ReleaseMonth uint      `json:"releaseMonth" gorm:"column:release_month; type:uint;"`
	ReleaseDay   uint      `json:"releaseDay" gorm:"column:release_day; type:uint;"`
}
