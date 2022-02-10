package models

type Country struct {
	Country string `json:"country" bson:"country"`
	Code    string `json:"code" bson:"code"`
}
