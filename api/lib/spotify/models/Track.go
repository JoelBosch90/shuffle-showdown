package spotify_models

type Track struct {
	PreviewUrl string `json:"preview_url"`
	IsPlayable bool   `json:"is_playable"`
	Type       string `json:"type"`
	Track      bool   `json:"track"`
	Album      `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	Id         string   `json:"id"`
	Name       string   `json:"name"`
}
