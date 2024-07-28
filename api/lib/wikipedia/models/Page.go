package wikipedia_models

type Page struct {
	PageId     int        `json:"pageid"`
	Title      string     `json:"title"`
	Revisions  []Revision `json:"revisions"`
	Categories []Category `json:"categories"`
}
