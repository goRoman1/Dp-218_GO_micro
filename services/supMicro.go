package services

import (
	"Dp218GO/models"
	"Dp218GO/repositories"
)

// ScooterInitService - structure for implementing ScooterInitService
type SupMicroService struct {
	supMicroRepo repositories.SupMicroRepoI
}

// NewScooterInitService - initialization of ScooterInitService
func NewSupMicroService(supMicroRepo repositories.SupMicroRepoI) *SupMicroService {
	return &SupMicroService{supMicroRepo}
}

// GetOwnersScooters - get all scooters that related to current user
func (si *SupMicroService) GetStations() (*models.StationList, error) {
	return si.supMicroRepo.GetStations()
}

// GetActiveStations - get all active system stations
func (si *SupMicroService) GetLocations() (*models.LocationList, error) {
	return si.supMicroRepo.GetLocations()
}

// AddStatusesToScooters - add statuses to scooter statuses
func (si *SupMicroService) AddStation(station *models.Station) error {
	return si.supMicroRepo.CreateStation(station)
}

func (si *SupMicroService) AddStationInLocation(station *models.Station, location *models.Location) error {
	return si.supMicroRepo.CreateStationInLocation(station, location)
}
