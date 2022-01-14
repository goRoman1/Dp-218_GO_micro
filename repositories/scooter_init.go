//go:generate mockgen -source=scooter_init.go -destination=../repositories/mocks/mock_scooter_init_repository.go -package=mock

package repositories

import "Dp-218_GO_micro/models"

// ScooterInitRepoI - interface for adding scooters to stations
type ScooterInitRepoI interface {
	GetOwnersScooters() (*models.SuppliersScooterList, error)
	GetActiveStations() (*models.StationList, error)
	AddStatusesToScooters(scooterIds []int, station models.Station) error
}
