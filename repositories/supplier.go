//go:generate mockgen -source=supplier.go -destination=../repositories/mocks/mock_supplier_repository.go -package=mock
package repositories

import (
	"Dp218GO/models"
)

// SupplierRepoI - interface for supplier repository
type SupplierRepoI interface {
	GetModels() (*models.ScooterModelDTOList, error)
	SelectModel(id int) (*models.ScooterModelDTO, error)
	AddModel(modelData *models.ScooterModelDTO) error
	EditPrice(modelData *models.ScooterModelDTO) error

	AddSuppliersScooter(modelId int, scooter string) error
	DeleteSuppliersScooter(id int) error
	ConvertToStruct(path string) []models.UploadedScooters
	InsertToDb(modelId int, scooters []models.UploadedScooters) error
}
