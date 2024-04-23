package spotify_models

type Album struct {
	IsPlayable           bool     `json:"is_playable"`
	Type                 string   `json:"type"`
	AlbumType            string   `json:"album_type"`
	Href                 string   `json:"href"`
	Id                   string   `json:"id"`
	Images               []Image  `json:"images"`
	Name                 string   `json:"name"`
	ReleaseDate          string   `json:"release_date"`
	ReleaseDatePrecision string   `json:"release_date_precision"`
	URI                  string   `json:"uri"`
	Artists              []Artist `json:"artists"`
	ExternalUrls         `json:"external_urls"`
	TotalTracks          int `json:"total_tracks"`
}
