package models

type PlayList struct {
	BaseEntity
	ID            string `json:"id" gorm:"type:string;primary_key;"`
	Name          string `json:"name" gorm:"type:string;"`
	LastSongAdded string `json:"lastSongAdded" gorm:"type:string;"`
	TracksTotal   uint   `json:"tracksTotal" gorm:"type:uint;"`
}
