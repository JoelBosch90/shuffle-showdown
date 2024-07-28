package wikipedia_models

type Page struct {
	PageId     int        `json:"pageid"`
	Title      string     `json:"title"`
	Categories []Category `json:"categories"`
	Links      []Link     `json:"links"`
}
