package models

type Query struct {
	Redirects []Redirect `json:"redirects"`
	Pages     []Page     `json:"pages"`
}
