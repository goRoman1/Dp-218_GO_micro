package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
)

type StationService struct {
	repoStation repositories.StationRepo
}

func NewStationService(repoStation repositories.StationRepo) *StationService {
	return &StationService{repoStation: repoStation}
}

func (db *StationService) GetAllStations() (*models.StationList, error) {
	return db.repoStation.GetAllStations()
}

func (db *StationService) AddStation(station *models.Station) error {
	return db.repoStation.AddStation(station)
}

func (db *StationService) GetStationById(stationId int) (models.Station, error) {
	return db.repoStation.GetStationById(stationId)
}

func (db *StationService) DeleteStation(stationId int) error {
	return db.repoStation.DeleteStation(stationId)
}

func (ser *StationService) ChangeStationBlockStatus(stationId int) error {
	station, err := ser.repoStation.GetStationById(stationId)
	if err != nil {
		return err
	}
	station.IsActive = !station.IsActive
	_, err = ser.repoStation.UpdateStation(stationId, station)
	return err
}

func (ser *StationService) UpdateStation(stationId int, stationData models.Station) (models.Station, error) {
	return ser.repoStation.UpdateStation(stationId, stationData)
}
