package models

type Game struct {
	OwnedEntity
	PlayListID uint `json:"playlistId" gorm:"type:uint;column:playlist_foreign_key;not null;"`
}
