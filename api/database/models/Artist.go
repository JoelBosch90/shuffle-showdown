package models

import "time"

type Artist struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Id        string    `json:"id" gorm:"column:id; type:string; primary_key;"`
	Name      string    `json:"name" gorm:"column:name; type:string;"`
}
