package models

import "time"

type Reception struct {
	ID       string    `json:"id"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
	DateTime time.Time `json:"DateTime"`
}

type FullReception struct {
	ID       string    `json:"id"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
	DateTime time.Time `json:"DateTime"`
	Products []Product `json:"Products"`
}
