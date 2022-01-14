package usecases

import "Dp-218_GO_micro/models"

type StationUsecasesRepo interface {
	GetAllStations() (*models.StationList, error)
	GetStationById(stationId int) (models.Station, error)
	AddStation(station *models.Station) error
	DeleteStation(stationId int) error
	ChangeStationBlockStatus(stationId int) error
}