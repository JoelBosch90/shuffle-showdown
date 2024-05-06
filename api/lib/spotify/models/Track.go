package spotify_models

type Track struct {
	PreviewUrl   string `json:"preview_url"`
	IsPlayable   bool   `json:"is_playable"`
	Explicit     bool   `json:"explicit"`
	Type         string `json:"type"`
	Episode      bool   `json:"episode"`
	Track        bool   `json:"track"`
	Album        `json:"album"`
	Artists      []Artist `json:"artists"`
	DiscNumber   int      `json:"disc_number"`
	TrackNumber  int      `json:"track_number"`
	DurationMs   int      `json:"duration_ms"`
	ExternalIds  `json:"external_ids"`
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Popularity   int    `json:"popularity"`
	URI          string `json:"uri"`
	IsLocal      bool   `json:"is_local"`
}
