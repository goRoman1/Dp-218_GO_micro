package models

import "time"

//ScooterDTO is a scooter model with custom parameters.
type ScooterDTO struct {
	ID            int     `json:"scooter_id"`
	ScooterModel  string  `json:"scooter_model"`
	MaxWeight     float64 `json:"max_weight"`
	BatteryRemain float64 `json:"battery_remain"`
	CanBeRent     bool    `json:"can_be_rent"`
}

//ScooterListDTO keeps a list of ScooterDTO.
type ScooterListDTO struct {
	Scooters []ScooterDTO `json:"scooters"`
}

//ScooterStatus keeps values of dynamic scooter parameters.
type ScooterStatus struct {
	Scooter       ScooterDTO `json:"scooter"`
	Location      Coordinate `json:"location"`
	BatteryRemain float64    `json:"battery_remain"`
	StationID     int        `json:"station_id"`
}

//ScooterStatusInRent keeps values which are important for the start and the end of the trip.
type ScooterStatusInRent struct {
	ID        int        `json:"id"`
	StationID int        `json:"station_id"`
	DateTime  time.Time  `json:"date_time"`
	Location  Coordinate `json:"location"`
}
