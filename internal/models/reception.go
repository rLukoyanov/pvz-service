package models

import "time"

type Reception struct {
	ID       string    `json:"id"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
	DateTime time.Time `json:"DateTime"`
}
