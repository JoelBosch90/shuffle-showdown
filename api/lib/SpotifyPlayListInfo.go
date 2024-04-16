package lib

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type Image struct {
	Height int    `json:"height"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
}

type Owner struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type AddedBy struct {
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type Artist struct {
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type ExternalIds struct {
	Isrc string `json:"isrc"`
}

type VideoThumbnail struct {
	Url string `json:"url"`
}

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

type Item struct {
	AddedAt        string `json:"added_at"`
	AddedBy        `json:"added_by"`
	IsLocal        bool   `json:"is_local"`
	PrimaryColor   string `json:"primary_color"`
	Track          `json:"track"`
	VideoThumbnail `json:"video_thumbnail"`
}

type SpotifyPlayListInfo struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  `json:"external_urls"`
	Followers     `json:"followers"`
	Href          string  `json:"href"`
	Id            string  `json:"id"`
	Images        []Image `json:"images"`
	Name          string  `json:"name"`
	Owner         `json:"owner"`
	PrimaryColor  string `json:"primary_color"`
	Public        bool   `json:"public"`
	SnapshotId    string `json:"snapshot_id"`
	Tracks        struct {
		Href     string `json:"href"`
		Items    []Item `json:"items"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}
