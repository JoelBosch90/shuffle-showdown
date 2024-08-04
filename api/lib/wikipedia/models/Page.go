package models

type Page struct {
	PageId     int        `json:"pageid"`
	Title      string     `json:"title"`
	Missing    bool       `json:"missing"`
	Categories []Category `json:"categories"`
	Links      []Link     `json:"links"`
}
