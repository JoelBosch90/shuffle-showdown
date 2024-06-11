package spotify_models

type Playlist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Tracks `json:"tracks"`
}
