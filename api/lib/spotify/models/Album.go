package spotify_models

type Album struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	ReleaseDate          string   `json:"release_date"`
	ReleaseDatePrecision string   `json:"release_date_precision"`
	Artists              []Artist `json:"artists"`
}
