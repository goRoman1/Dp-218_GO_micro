package main

import (
	"OrderService/config"
	"OrderService/proto"
	"OrderService/repository"
	"OrderService/service"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	log.Println("Starting order microservice")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PG_HOST,
		config.PG_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "order_micro", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepo(db)
	service := service.NewOrderService(orderRepo)

	listener, err := net.Listen("tcp", net.JoinHostPort("", os.Getenv("ORDER_GRPC_PORT")))
	if err != nil {
		log.Panicf("%s: failed to listen on port - %v", "order_micro", err)
	}

	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, service)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Panicf("%s: failed to start grpc - %v", "order_micro", err)
	}
}
