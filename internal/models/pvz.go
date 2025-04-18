package models

import "time"

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type FullPVZ struct {
	ID               string                   `json:"id"`
	RegistrationDate time.Time                `json:"registrationDate"`
	City             string                   `json:"city"`
	Receptions       map[string]FullReception `json:"receptions"`
}
