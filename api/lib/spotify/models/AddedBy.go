package spotify_models

type AddedBy struct {
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}
