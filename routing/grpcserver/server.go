package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

//NewGrpcServer creates a new gRPC server on port 8080.
func NewGrpcServer() *grpc.Server{
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		fmt.Println("grpc server started: 8000")
		log.Fatal(grpcServer.Serve(listener))
	}()
	return grpcServer
}