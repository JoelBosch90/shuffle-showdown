package wikipedia_models

type Query struct {
	Pages []Page `json:"pages"`
}
