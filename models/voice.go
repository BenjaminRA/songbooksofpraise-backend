package models

type Voice struct {
	Voice    string  `json:"voice"`
	File     string  `json:"file"`
	Duration float64 `json:"duration"`
}
