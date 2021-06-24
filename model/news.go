package model

type News struct {
	ID          string   `json:"ID"`
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	Tag         []string `json:"Tag"`
	Status      string   `json:"Status"`
}
