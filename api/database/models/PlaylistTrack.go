package models

type PlaylistTrack struct {
	PlaylistId string `json:"playlistId" gorm:"type:string; uniqueIndex:uniquelink;not null;"`
	TrackId    string `json:"trackId" gorm:"type:string; uniqueIndex:uniquelink;not null;"`
}
