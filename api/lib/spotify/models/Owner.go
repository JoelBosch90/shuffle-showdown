package spotify_models

type Owner struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls `json:"external_urls"`
	Href         string `json:"href"`
	Id           string `json:"id"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}
