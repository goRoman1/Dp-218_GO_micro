package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
)

// ScooterInitService - structure for implementing ScooterInitService
type ScooterInitService struct {
	scooterInitRepo repositories.ScooterInitRepoI
}

// NewScooterInitService - initialization of ScooterInitService
func NewScooterInitService(scooterInitRepo repositories.ScooterInitRepoI) *ScooterInitService {
	return &ScooterInitService{scooterInitRepo}
}

// GetOwnersScooters - get all scooters that related to current user
func (si *ScooterInitService) GetOwnersScooters() (*models.SuppliersScooterList, error) {
	return si.scooterInitRepo.GetOwnersScooters()
}

// GetActiveStations - get all active system stations
func (si *ScooterInitService) GetActiveStations() (*models.StationList, error) {
	return si.scooterInitRepo.GetActiveStations()
}

// AddStatusesToScooters - add statuses to scooter statuses
func (si *ScooterInitService) AddStatusesToScooters(scooterIds []int, station models.Station) error {
	return si.scooterInitRepo.AddStatusesToScooters(scooterIds, station)
}

// ConvertForTemplateStruct - create struct for template rendering
func (si *ScooterInitService) ConvertForTemplateStruct() *models.ScootersStationsAllocation {
	list := &models.ScootersStationsAllocation{}

	scooters, err := si.scooterInitRepo.GetOwnersScooters()
	if err != nil {
		return nil
	}
	stations, err := si.scooterInitRepo.GetActiveStations()
	if err != nil {
		return nil
	}

	list.SuppliersScooterList = *scooters
	list.StationList = *stations

	return list
}
