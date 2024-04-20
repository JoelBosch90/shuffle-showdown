package models

type Artist struct {
	BaseEntity
	ID   string `json:"id" gorm:"type:string;primary_key;"`
	Name string `json:"name" gorm:"type:string;"`
}
