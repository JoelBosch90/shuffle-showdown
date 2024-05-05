package models

type TrackArtist struct {
	TrackId  string `json:"trackId" gorm:"type:string; uniqueIndex:uniquelink;not null;"`
	ArtistId string `json:"artistId" gorm:"type:string; uniqueIndex:uniquelink;not null;"`
}
