package spotify_models

type Item struct {
	AddedAt        string `json:"added_at"`
	AddedBy        `json:"added_by"`
	IsLocal        bool   `json:"is_local"`
	PrimaryColor   string `json:"primary_color"`
	Track          `json:"track"`
	VideoThumbnail `json:"video_thumbnail"`
}
