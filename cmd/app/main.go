package main

import (
	"Dp-218_GO_micro/configs"
	"Dp-218_GO_micro/protos"
	"Dp-218_GO_micro/repositories/postgres"
	"Dp-218_GO_micro/routing"
	"Dp-218_GO_micro/routing/grpcserver"
	"Dp-218_GO_micro/routing/httpserver"
	"Dp-218_GO_micro/services"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
)

func main() {

	var connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_DB)

	db, err := postgres.NewConnection(connectionString)
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer db.CloseDB()

	err = doMigrate(connectionString)
	if err != nil {
		log.Printf("app - Run - Migration issues: %v\n", err)
	}

	var userRoleRepoDB = postgres.NewUserRepoDB(db)
	var userService = services.NewUserService(userRoleRepoDB, userRoleRepoDB)

	var accRepoDB = postgres.NewAccountRepoDB(userRoleRepoDB, db)
	var clock = services.NewClock()
	var accService = services.NewAccountService(accRepoDB, accRepoDB, accRepoDB, clock)
	var stationRepoDB = postgres.NewStationRepoDB(db)
	var stationService = services.NewStationService(stationRepoDB)

	var scooterRepo = postgres.NewScooterRepoDB(db)
	var grpcScooterService = services.NewGrpcScooterService(scooterRepo, stationService)
	var scooterService = services.NewScooterService(scooterRepo)

	var supplierRepoDB = postgres.NewSupplierRepoDB(db)
	var supplierService = services.NewSupplierService(supplierRepoDB)

	problemGRPCServer := net.JoinHostPort(configs.PROBLEMS_SERVICE, configs.PROBLEMS_GRPC_PORT)
	problemCred, err := credentials.NewClientTLSFromFile(configs.PROBLEMS_CERTIFICATE, "")
	if err != nil {
		log.Panicf("%s: unable to get TLS certificate - %v", problemGRPCServer, err)
	}
	problemConnection, err := grpc.Dial(problemGRPCServer, grpc.WithTransportCredentials(problemCred))
	if err != nil {
		log.Panicf("%s: unable to set grpc connection - %v", problemGRPCServer, err)
	}
	defer problemConnection.Close()
	var problemService = services.NewProblemService(problemConnection, userService)

	var orderRepoDB = postgres.NewOrderRepoDB(db)
	var orderService = services.NewOrderService(orderRepoDB)

	var scootersInitRepoDb = postgres.NewScooterInitRepoDB(db)
	var scootersInitService = services.NewScooterInitService(scootersInitRepoDb)

	sessStore := sessions.NewCookieStore([]byte(configs.SESSION_SECRET))
	authService := services.NewAuthService(userRoleRepoDB, sessStore)

	custService := services.NewCustomerService(stationRepoDB)

	/////////////////
	/*
		supplierMicroGRPCServer := net.JoinHostPort(configs.SUPPLIER_MICRO_SERVICE, configs.SUPPLIER_MICRO_GRPC_PORT)
		supplierMicroCred, err := credentials.NewClientTLSFromFile(configs.SUPPLIER_MICRO_CERTIFICATE, "")
		if err != nil {
			log.Panicf("%s: unable to get TLS certificate - %v", supplierMicroGRPCServer, err)
		}
		supplierMicroConnection, err := grpc.Dial(supplierMicroGRPCServer, grpc.WithTransportCredentials(supplierMicroCred))
		if err != nil {
			log.Panicf("%s: unable to set grpc connection - %v", supplierMicroGRPCServer, err)
		}
		defer supplierMicroConnection.Close()
		var supplierMicroService = services.NewSupplierMicroService(supplierMicroConnection, userService)

	*/
	/////////////////

	handler := routing.NewRouter()
	routing.AddAuthHandler(handler, authService)
	routing.AddCustomerHandler(handler, custService)
	routing.AddUserHandler(handler, userService)
	routing.AddStationHandler(handler, stationService)
	routing.AddAccountHandler(handler, accService)
	routing.AddScooterHandler(handler, scooterService)
	routing.AddProblemHandler(handler, problemService)
	//	routing.AddSupplierMicroHandler(handler, supplierMicroService)
	routing.AddGrpcScooterHandler(handler, grpcScooterService)
	routing.AddOrderHandler(handler, orderService)
	routing.AddSupplierHandler(handler, supplierService)
	routing.AddScooterInitHandler(handler, scootersInitService)
	httpServer := httpserver.New(handler, httpserver.Port(configs.HTTP_PORT))
	handler.HandleFunc("/scooter", httpServer.ScooterHandler)

	//utils.CheckKafka() //TODO: delete after checking

	grpcServer := grpcserver.NewGrpcServer()
	protos.RegisterScooterServiceServer(grpcServer, httpServer)
	http.ListenAndServe(":8080", handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatalf("app - Run - httpServer.Notify: %v", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Fatalf("app - Run - httpServer.Shutdown: %v", err)
	}
}

func doMigrate(connStr string) error {
	migr, err := migrate.New("file://"+configs.MIGRATIONS_PATH, connStr+"?sslmode=disable")
	if err != nil {
		return err
	}

	if configs.MIGRATE_VERSION_FORCE > 0 {
		migr.Force(configs.MIGRATE_VERSION_FORCE)
	}

	if configs.MIGRATE_DOWN {
		migr.Down()
	}

	return migr.Up()
}
