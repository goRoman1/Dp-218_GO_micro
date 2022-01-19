package services

import (
	"Dp218GO/models"
	"context"
	"google.golang.org/grpc"
	proto "supplier.micro/proto"
)

// SupplierMicroService - structure for implementing user problem service
type SupplierMicroService struct {
	microservice proto.SupplierMicroServiceClient
	userService  *UserService
}

// NewSupplierMicroService - initialization of ProblemService
func NewSupplierMicroService(grpcConn grpc.ClientConnInterface, userServ *UserService) *SupplierMicroService {
	return &SupplierMicroService{
		microservice: proto.NewSupplierMicroServiceClient(grpcConn),
		userService:  userServ,
	}
}

func (supserv *SupplierMicroService) unmarshallLocations(locationsTypeGRPC *proto.Location) models.Location {
	return models.Location{
		ID:        int(locationsTypeGRPC.Id),
		Longitude: float64(locationsTypeGRPC.Longitude),
		Latitude:  float64(locationsTypeGRPC.Latitude),
		Label:     locationsTypeGRPC.Label,
	}
}

func (supserv *SupplierMicroService) marshallStation(station *models.Station) *proto.Station {
	return &proto.Station{
		Id:        int32(station.ID),
		Name:      station.Name,
		IsActive:  station.IsActive,
		Latitude:  float32(station.Latitude),
		Longitude: float32(station.Longitude),
	}
}

func (supserv *SupplierMicroService) AddNewStation(station *models.Station) error {
	stationToAdd := &proto.Station{
		Id:        int32(station.ID),
		Name:      station.Name,
		IsActive:  station.IsActive,
		Latitude:  float32(station.Latitude),
		Longitude: float32(station.Longitude),
	}
	_, err := supserv.microservice.CreateStation(context.Background(), stationToAdd)
	return err
}

func (supserv *SupplierMicroService) GetAllLocations() ([]models.Location, error) {
	request := &proto.Request{}
	response, err := supserv.microservice.GetLocations(context.Background(), request)
	var locationList []models.Location
	if err != nil {
		return locationList, err
	}
	for _, val := range response.Locations {
		locationList = append(locationList, supserv.unmarshallLocations(val))
	}
	return locationList, err
}

// AddProblemSolution - make solution record for given problem (by ID)
func (supserv *SupplierMicroService) CreateStationInLocation(location *models.Location, station *models.Station) error {
	request := &proto.StationLocation{
		Location: &proto.Location{
			Longitude: float32(location.Longitude),
			Latitude:  float32(location.Latitude),
		},
		ScooterStation: supserv.marshallStation(station),
	}
	_, err := supserv.microservice.CreateStationInLocation(context.Background(), request)
	return err
}
