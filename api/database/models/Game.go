package models

type Game struct {
	OwnedEntity
	PlayListID     uint `json:"playlistId" gorm:"type:uint;column:playlist_foreign_key;not null;"`
	SongsToWin     uint `json:"songsToWin" gorm:"type:uint;"`
	TitleRequired  bool `json:"titleRequired" gorm:"type:bool;"`
	ArtistRequired bool `json:"artistRequired" gorm:"type:bool;"`
	Configured     bool `json:"configured" gorm:"type:bool;"`
}
