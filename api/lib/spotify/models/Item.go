package spotify_models

type Item struct {
	AddedAt string `json:"added_at"`
	Track   `json:"track"`
}
