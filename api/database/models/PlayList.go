package models

type PlayList struct {
	BaseEntity
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
