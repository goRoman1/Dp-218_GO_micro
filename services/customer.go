package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"math"
)

// CustomerService takes repostation interface
// serves as main service for customer interaction with system
type CustomerService struct {
	repoStation repositories.StationRepo
}

// NewCustomerService returns new customerservice
func NewCustomerService(repo repositories.StationRepo) *CustomerService {
	return &CustomerService{
		repoStation: repo,
	}
}

// ListStations returns station list from db or error if failed
func (cs *CustomerService) ListStations() (*models.StationList, error) {
	return cs.repoStation.GetAllStations()
}

// ShowStation returns station from db by id or error if failed
func (cs *CustomerService) ShowStation(id int) (*models.Station, error) {
	station, err := cs.repoStation.GetStationById(id)
	if err != nil {
		return nil, err
	}
	return &station, nil
}

// ShowNearestStation takes user location and returns nearest station or error if failed
func (cs *CustomerService) ShowNearestStation(x, y float64) (*models.Station, error) {

	stations, err := cs.repoStation.GetAllStations()
	if err != nil {
		return nil, err
	}

	nearest := calcNearest(x, y, stations.Station)
	return nearest, nil
}

func calcNearest(x, y float64, sts []models.Station) *models.Station {

	min := math.MaxFloat64
	var nearest models.Station

	for _, v := range sts {
		dis := calcDistance(x, y, v.Latitude, v.Longitude)
		if dis < min {
			min = dis
			nearest = v
		}
	}
	return &nearest
}

func calcDistance(x1, y1, x2, y2 float64) float64 {
	z := math.Pow(math.Abs(x1-x2), 2) + math.Pow(math.Abs(y1-y2), 2)
	return math.Sqrt(z)
}
