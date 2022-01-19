package repositories

import (
	"Dp218GO/models"
)

type StationRepo interface {
	GetAllStations() (*models.StationList, error)
	GetStationById(stationId int) (models.Station, error)
	AddStation(station *models.Station) error
	DeleteStation(stationId int) error
	UpdateStation(stationId int, stationData models.Station) (models.Station, error)
}
