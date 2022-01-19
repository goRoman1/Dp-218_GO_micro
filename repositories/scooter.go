//go:generate mockgen -source=scooter.go -destination=../repositories/mock/mock_scooter.go -package=mock
package repositories

import "Dp218GO/models"

//ScooterRepo the interface which implemented by functions which connect to the database.
type ScooterRepo interface {
	GetAllScooters() (*models.ScooterListDTO, error)
	GetAllScootersByStationID(stationID int) (*models.ScooterListDTO, error)
	GetScooterById(scooterId int) (models.ScooterDTO, error)
	GetScooterStatus(scooterID int) (models.ScooterStatus, error)
	SendCurrentStatus(id, stationID int, lat, lon, battery float64) error
	CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error)
}
