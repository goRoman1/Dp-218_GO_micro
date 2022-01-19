package services

import (
	"Dp218GO/models"
	"Dp218GO/protos"
	"Dp218GO/repositories"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

const (
	step          = 0.0001
	dischargeStep = 0.1
	interval      = 450
)

//GrpcScooterService is a service which responsible for gRPC scooter.
type GrpcScooterService struct {
	repositories.ScooterRepo
	*StationService
}

//GrpcScooterClient is a struct with parameters which will be translated by the gRPC connection.
type GrpcScooterClient struct {
	ID            uint64
	coordinate    models.Coordinate
	batteryRemain float64
	stream        protos.ScooterService_ReceiveClient
}

//NewGrpcScooterService creates a new GrpcScooterService.
func NewGrpcScooterService(repoScooter repositories.ScooterRepo, stationService *StationService) *GrpcScooterService {
	return &GrpcScooterService{
		repoScooter,
		stationService,
	}
}

//NewGrpcScooterClient creates a new GrpcScooterClient with given parameters.
func NewGrpcScooterClient(id uint64, coordinate models.Coordinate, battery float64,
	stream protos.ScooterService_ReceiveClient) *GrpcScooterClient {
	return &GrpcScooterClient{
		ID:            id,
		coordinate:    coordinate,
		batteryRemain: battery,
		stream:        stream,
	}
}

//InitAndRun the main function of scooter's trip. It analyzes the scooter parameters from database by its ID.
//If they satisfy the conditions, function creates connection to the gRPC server, creates gRPC client,
//calls 'run' function which moves the scooter to the destination point.
//After finished moves it sends the current scooter status to the database.
func (gss *GrpcScooterService) InitAndRun(scooterID int, chosenStationID int) error {
	scooter, err := gss.GetScooterById(scooterID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	scooterStatus, err := gss.GetScooterStatus(scooterID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if scooter.CanBeRent {
		var coordinate models.Coordinate
		station, err := gss.GetStationById(chosenStationID)
		if err != nil {
			return err
		}
		coordinate.Latitude = station.Latitude
		coordinate.Longitude = station.Longitude

		conn, err := grpc.DialContext(context.Background(), ":8000", grpc.WithInsecure())

		if err != nil {
			panic(err)
		}
		defer conn.Close()

		sClient := protos.NewScooterServiceClient(conn)
		stream, err := sClient.Receive(context.Background())
		if err != nil {
			panic(err)
		}

		client := NewGrpcScooterClient(uint64(scooterID),
			scooterStatus.Location, scooter.BatteryRemain, stream)
		err = client.run(coordinate)
		if err != nil {
			fmt.Println(err)
		}

		err = gss.SendCurrentStatus(int(client.ID), chosenStationID, client.coordinate.Latitude,
			client.coordinate.Longitude,
			client.batteryRemain)
		if err != nil {
			fmt.Println(err)
		}

		if client.batteryRemain <= 0 {
			err = fmt.Errorf("scooter battery discharged. Trip is over")
			return err
		}
		return nil
	}

	err = fmt.Errorf("you can't use this scooter. Choose another one")
	fmt.Println(err.Error())
	return err
}

//grpcScooterMessage sends the message be gRPC stream in a format which defined in the *proto file.
func (s *GrpcScooterClient) grpcScooterMessage() {
	intPol := time.Duration(interval) * time.Millisecond

	fmt.Println("executing run in client")
	msg := &protos.ClientMessage{
		Id:        s.ID,
		Latitude:  s.coordinate.Latitude,
		Longitude: s.coordinate.Longitude,
	}
	err := s.stream.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(intPol)
}

//run is responsible for scooter's movements from his current position to the destination point.
//Run also is responsible for scooter's discharge. Every step battery charge decrease by the constant discharge value.
func (s *GrpcScooterClient) run(station models.Coordinate) error {

	switch {
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.
			coordinate.Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude+step, s.coordinate.Longitude+step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude-step, s.coordinate.Longitude+step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude-step, s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude+step, s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude <= station.Latitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.
			batteryRemain > 0; s.coordinate.Latitude, s.batteryRemain = s.coordinate.Latitude+step, s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.
			batteryRemain > 0; s.coordinate.Latitude, s.batteryRemain = s.coordinate.Latitude-step, s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.Longitude, s.batteryRemain = s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.coordinate.Longitude, s.batteryRemain = s.coordinate.Longitude+step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
	default:
		return fmt.Errorf("error happened")
	}
	return nil
}
