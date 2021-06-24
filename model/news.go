package model

type News struct {
	NID        string `json:"NID"`
	Title      string `json:"Title"`
	Desciption string `json:"Description"`
	Tag        []Tags `json:"Tags"`
}

type Tags struct {
	TID string `json:"TID"`
	Tag string `json:"Tag"`
}
