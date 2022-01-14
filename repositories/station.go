package repositories

import (
	"Dp-218_GO_micro/models"
)

type StationRepo interface {
	GetAllStations() (*models.StationList, error)
	GetStationById(stationId int) (models.Station, error)
	AddStation(station *models.Station) error
	DeleteStation(stationId int) error
	UpdateStation(stationId int, stationData models.Station) (models.Station, error)
}
