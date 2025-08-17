package models

type Verse struct {
	OrderNumber int    `json:"order_number"`
	Text        string `json:"text"`
	Chorus      bool   `json:"chorus"`
}
