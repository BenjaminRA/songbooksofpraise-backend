package models

type Verse struct {
	OrderNumber int    `json:"order_number" bson:"order_number"`
	Text        string `json:"text" bson:"text"`
	Chorus      bool   `json:"chorus" bson:"chorus"`
}
