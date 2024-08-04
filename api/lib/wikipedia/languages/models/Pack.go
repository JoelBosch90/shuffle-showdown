package models

type Pack struct {
	Language           Language           `json:"language"`
	ApiUrl             string             `json:"apiUrl"`
	CategoryParsers    CategoryParsers    `json:"categoryParsers"`
	RedirectFormatters RedirectFormatters `json:"redirectFormatters"`
}
