package spotify_models

type PlayList struct {
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
	Tracks        `json:"tracks"`
	Type          string `json:"type"`
	URI           string `json:"uri"`
}
