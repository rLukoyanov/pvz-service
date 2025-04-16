package models

type Reception struct {
	ID       string `json:"id"`
	PvzId    string `json:"pvzId"`
	Status   string `json:"status"`
	DateTime string `json:"DateTime"`
}
