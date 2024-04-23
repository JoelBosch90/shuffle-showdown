package spotify_models

type Artist struct {
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}
