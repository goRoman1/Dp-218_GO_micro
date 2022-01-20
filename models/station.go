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

type Location struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Label     string  `json:"label"`
}

type LocationList struct {
	Location []Location `json:"location"`
}

type StationLocation struct {
	StationList  []Station  `json:"station_list"`
	LocationList []Location `json:"location_list"`
}
