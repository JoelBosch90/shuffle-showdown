package models

type Track struct {
	BaseEntity
	ID           string `json:"id" gorm:"type:string;primary_key;"`
	Name         string `json:"name" gorm:"type:string;"`
	ReleaseYear  uint   `json:"releaseYear" gorm:"type:uint;"`
	ReleaseMonth uint   `json:"releaseMonth" gorm:"type:uint;"`
	ReleaseDay   uint   `json:"releaseDay" gorm:"type:uint;"`
}
