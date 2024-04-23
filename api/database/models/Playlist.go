package models

type Playlist struct {
	BaseEntity
	ID            string `json:"id" gorm:"column:id; type:string; primary_key;"`
	Name          string `json:"name" gorm:"column:name; type:string;"`
	LastSongAdded string `json:"lastSongAdded" gorm:"column:last_song_added; type:string;"`
	TracksTotal   uint   `json:"tracksTotal" gorm:"column:tracks_total; type:uint;"`
}
