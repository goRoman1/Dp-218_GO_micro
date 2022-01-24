package service

import (
	"ScooterServer/config"
	"ScooterServer/proto"
	"ScooterServer/repository"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	step          = 0.0001
	dischargeStep = 0.1
	interval      = 450
)

type Location struct {
	Latitude  float64
	Longitude float64
}

//ScooterService is a service which responsible for gRPC scooter.
type ScooterService struct {
	Repo  *repository.ScooterRepo
	Order proto.OrderServiceClient
}

//ScooterClient is a struct with parameters which will be translated by the gRPC connection.
type ScooterClient struct {
	ID            uint64
	Latitude      float64
	Longitude     float64
	BatteryRemain float64
	Stream        proto.ScooterService_ReceiveClient
}

//NewScooterService creates a new GrpcScooterService.
func NewScooterService(repoScooter *repository.ScooterRepo, order proto.OrderServiceClient) *ScooterService {
	return &ScooterService{
		Repo:  repoScooter,
		Order: order,
	}
}

//NewScooterClient creates a new GrpcScooterClient with given parameters.
func NewScooterClient(id uint64, latitude, longitude, battery float64,
	stream proto.ScooterService_ReceiveClient) *ScooterClient {
	return &ScooterClient{
		ID:            id,
		Latitude:      latitude,
		Longitude:     longitude,
		BatteryRemain: battery,
		Stream:        stream,
	}
}

//InitAndRun the main function of scooter's trip. It analyzes the scooter parameters from database by its ID.
//If they satisfy the conditions, function creates connection to the gRPC server, creates gRPC client,
//calls 'run' function which moves the scooter to the destination point.
//After finished moves it sends the current scooter status to the database.
func (gss *ScooterService) InitAndRun(ctx context.Context, id *proto.ScooterID, stationID *proto.StationID) error {
	scooter, err := gss.GetScooterById(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	scooterStatus, err := gss.GetScooterStatus(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if scooter.CanBeRent {
		var coordinate Location
		station, err := gss.GetStationById(ctx, stationID)
		if err != nil {
			return err
		}
		coordinate.Latitude = station.Latitude
		coordinate.Longitude = station.Longitude

		conn, err := grpc.DialContext(ctx, net.JoinHostPort("", config.GRPC_PORT), grpc.WithInsecure())

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		sClient := proto.NewScooterServiceClient(conn)
		stream, err := sClient.Receive(ctx)
		if err != nil {
			log.Fatal(err)
		}

		client := NewScooterClient(uint64(id.Id),
			scooterStatus.Latitude, scooterStatus.Longitude, scooter.BatteryRemain, stream)
		err = client.run(coordinate)
		if err != nil {
			fmt.Println(err)
		}

		sendStatus := &proto.SendStatus{
			ScooterID: client.ID, StationID: stationID.Id,
			Latitude: client.Latitude, Longitude: client.Longitude, BatteryRemain: client.BatteryRemain}

		_, err = gss.SendCurrentStatus(ctx, sendStatus)
		if err != nil {
			fmt.Println(err)
		}

		if client.BatteryRemain <= 0 {
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
func (s *ScooterClient) grpcScooterMessage() {
	intPol := time.Duration(interval) * time.Millisecond

	fmt.Println("executing run in client")
	msg := &proto.ClientMessage{
		Id:        s.ID,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}
	err := s.Stream.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(intPol)
}

//run is responsible for scooter's movements from his current position to the destination point.
//Run also is responsible for scooter's discharge. Every step battery charge decrease by the constant discharge value.
func (s *ScooterClient) run(station Location) error {

	switch {
	case s.Latitude <= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.
			Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude+step, s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude >= station.Latitude && s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude-step, s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude >= station.Latitude && s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude-step, s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Latitude <= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude+step, s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Latitude <= station.Latitude:
		for ; s.Latitude <= station.Latitude && s.
			BatteryRemain > 0; s.Latitude, s.BatteryRemain = s.Latitude+step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude:
		for ; s.Latitude >= station.Latitude && s.
			BatteryRemain > 0; s.Latitude, s.BatteryRemain = s.Latitude-step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Longitude >= station.Longitude:
		for ; s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Longitude, s.BatteryRemain = s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.Longitude <= station.Longitude:
		for ; s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.Longitude, s.BatteryRemain = s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
	default:
		return fmt.Errorf("error happened")
	}
	return nil
}

//GetAllScooters gives the access to the ScooterRepo.GetAllScooters function.
func (gss *ScooterService) GetAllScooters(ctx context.Context, request *proto.Request) (*proto.ScooterList, error) {
	return gss.Repo.GetAllScooters(ctx, request)
}

func (gss *ScooterService) GetAllScootersByStationID(ctx context.Context, id *proto.StationID) (*proto.ScooterList,
	error) {
	return gss.Repo.GetAllScootersByStationID(ctx, id)
}

func (gss *ScooterService) GetAllStations(ctx context.Context, request *proto.Request) (*proto.StationList,
	error) {
	return gss.Repo.GetAllStations(ctx, request)
}

//GetScooterById gives the access to the ScooterRepo.GetScooterById function.
func (gss *ScooterService) GetScooterById(ctx context.Context, id *proto.ScooterID) (*proto.Scooter, error) {
	return gss.Repo.GetScooterById(ctx, id)
}

func (gss *ScooterService) GetStationById(ctx context.Context, id *proto.StationID) (*proto.Station, error) {
	return gss.Repo.GetStationById(ctx, id)
}

//GetScooterStatus gives the access to the ScooterRepo.GetScooterStatus function.
func (gss *ScooterService) GetScooterStatus(ctx context.Context, status *proto.ScooterID) (*proto.ScooterStatus, error) {
	return gss.Repo.GetScooterStatus(ctx, status)
}

//SendCurrentStatus gives the access to the ScooterRepo.SendCurrentStatus function.
func (gss *ScooterService) SendCurrentStatus(ctx context.Context, status *proto.SendStatus) (*proto.Response, error) {
	return gss.Repo.SendCurrentStatus(ctx, status)
}

//CreateScooterStatusInRent gives the access to the ScooterRepo.CreateScooterStatusInRent function.
func (gss *ScooterService) CreateScooterStatusInRent(ctx context.Context, id *proto.ScooterID) (*proto.ScooterStatusInRent,
	error) {
	return gss.Repo.CreateScooterStatusInRent(ctx, id)
}
