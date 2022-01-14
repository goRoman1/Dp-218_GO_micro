package models

type Station struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	IsActive  bool    `json:"is_active"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type StationList struct {
	Station []Station `json:"station"`
}
