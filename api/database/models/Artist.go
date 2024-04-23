package models

type Artist struct {
	BaseEntity
	ID   string `json:"id" gorm:"column:id; type:string; primary_key;"`
	Name string `json:"name" gorm:"column:name; type:string;"`
}
