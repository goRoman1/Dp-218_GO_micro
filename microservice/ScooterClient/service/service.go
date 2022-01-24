package service

import (
	"ScooterClient/model"
	"ScooterClient/proto"
	"fmt"
	"time"
)

const (
	step          = 0.0001
	dischargeStep = 0.1
	interval      = 450
)

//ScooterClient is a struct with parameters which will be translated by the gRPC connection.
type ScooterClient struct {
	ID            uint64
	Latitude      float64
	Longitude     float64
	BatteryRemain float64
	Stream        proto.ScooterService_RegisterClient
}

//NewScooterClient creates a new GrpcScooterClient with given parameters.
func NewScooterClient(id uint64, latitude, longitude, battery float64,
	stream proto.ScooterService_RegisterClient) *ScooterClient {
	return &ScooterClient{
		ID:            id,
		Latitude:      latitude,
		Longitude:     longitude,
		BatteryRemain: battery,
		Stream:        stream,
	}
}

//GrpcScooterMessage sends the message be gRPC stream in a format which defined in the *proto file.
func (s *ScooterClient) GrpcScooterMessage() {
	intPol := time.Duration(interval) * time.Millisecond

	fmt.Println("executing run in client")
	msg := proto.ClientMessage{
		Id:        s.ID,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}

	fmt.Printf("Send to server this message: %v\n", &msg)
	err := s.Stream.Send(&msg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(intPol)
}

//Run is responsible for scooter's movements from his current position to the destination point.
//Run also is responsible for scooter's discharge. Every step battery charge decrease by the constant discharge value.
func (s *ScooterClient) Run(station model.Location) error {

	switch {
	case s.Latitude <= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.
			Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude+step, s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude >= station.Latitude && s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude-step, s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude >= station.Latitude && s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude-step, s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Latitude <= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Latitude,
			s.Longitude, s.BatteryRemain = s.Latitude+step, s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Latitude <= station.Latitude:
		for ; s.Latitude <= station.Latitude && s.
			BatteryRemain > 0; s.Latitude, s.BatteryRemain = s.Latitude+step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Latitude >= station.Latitude:
		for ; s.Latitude >= station.Latitude && s.
			BatteryRemain > 0; s.Latitude, s.BatteryRemain = s.Latitude-step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Longitude >= station.Longitude:
		for ; s.Longitude >= station.Longitude && s.
			BatteryRemain > 0; s.Longitude, s.BatteryRemain = s.Longitude-step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
		fallthrough
	case s.Longitude <= station.Longitude:
		for ; s.Longitude <= station.Longitude && s.
			BatteryRemain > 0; s.Longitude, s.BatteryRemain = s.Longitude+step,
			s.BatteryRemain-dischargeStep {
			s.GrpcScooterMessage()
		}
	default:
		return fmt.Errorf("error happened")
	}
	s.ID = 0
	return nil
}
