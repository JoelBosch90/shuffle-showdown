package models

type Response struct {
	BatchComplete bool  `json:"batchcomplete"`
	Query         Query `json:"query"`
}
