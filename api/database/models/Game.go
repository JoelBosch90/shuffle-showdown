package models

type Game struct {
	OwnedEntity
	PlayListId     string `json:"playListId" gorm:"type:string;column:playlist_foreign_key;not null;"`
	SongsToWin     uint   `json:"songsToWin" gorm:"type:uint;"`
	TitleRequired  bool   `json:"titleRequired" gorm:"type:bool;"`
	ArtistRequired bool   `json:"artistRequired" gorm:"type:bool;"`
	Configured     bool   `json:"configured" gorm:"type:bool;"`
}
