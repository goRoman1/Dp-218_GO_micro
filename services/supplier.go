package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
)

// SupplierService - structure for implementing supplier service
type SupplierService struct {
	SupplierRepo repositories.SupplierRepoI
}

// NewSupplierService - initialization of SupplierService
func NewSupplierService(SupplierRepo repositories.SupplierRepoI) *SupplierService {
	return &SupplierService{
		SupplierRepo: SupplierRepo,
	}
}

//AddSuppliersScooter - adds suppliers scooter to scooter model
func (s *SupplierService) AddSuppliersScooter(modelId int, scooter string) error {
	return s.SupplierRepo.AddSuppliersScooter(modelId, scooter)
}

// DeleteSuppliersScooter - delete suppliers scooter
func (s *SupplierService) DeleteSuppliersScooter(id int) error {
	return s.SupplierRepo.DeleteSuppliersScooter(id)
}

// InsertScootersToDb - adding scooters from .csv to Db
func (s *SupplierService) InsertScootersToDb(modelId int, path string) {
	scooterUploaded := s.SupplierRepo.ConvertToStruct(path)
	s.SupplierRepo.InsertToDb(modelId, scooterUploaded)
}

// GetModels - getting all scooter models
func (s *SupplierService) GetModels() (*models.ScooterModelDTOList, error) {
	return s.SupplierRepo.GetModels()
}

// SelectModel - get selected model data
func (s *SupplierService) SelectModel(id int) (*models.ScooterModelDTO, error) {
	return s.SupplierRepo.SelectModel(id)
}

// AddModel - adding of model to Db
func (s *SupplierService) AddModel(modelData *models.ScooterModelDTO) error {
	return s.SupplierRepo.AddModel(modelData)
}

// ChangePrice - change payment price that related to scooter model
func (s *SupplierService) ChangePrice(modelData *models.ScooterModelDTO) error {
	return s.SupplierRepo.EditPrice(modelData)
}
