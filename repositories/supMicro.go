package repositories

import "Dp218GO/models"

// ScooterInitRepoI - interface for adding scooters to stations
type SupMicroRepoI interface {
	GetStations() (*models.StationList, error)
	GetLocations() (*models.LocationList, error)
	CreateStationInLocation(station *models.Station, location *models.Location) error
	CreateStation(station *models.Station) error
}
