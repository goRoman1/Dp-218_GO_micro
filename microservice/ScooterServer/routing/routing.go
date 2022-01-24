package routing

import (
	"ScooterServer/config"
	"ScooterServer/proto"
	"ScooterServer/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

var (
	scooterIDKey = "scooterId"
	stationIDKey = "stationId"
)

var chosenScooterID, chosenStationID int

type combineForTemplate struct {
	*proto.ScooterList
	*proto.StationList
}

type Routing interface {
	getAllScooters(w http.ResponseWriter, r *http.Request)
	getScooterById(w http.ResponseWriter, r *http.Request)
	startScooterTrip(w http.ResponseWriter, r *http.Request)
	showTripPage(w http.ResponseWriter, r *http.Request)
	chooseScooter(w http.ResponseWriter, r *http.Request)
	chooseStation(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	scooterService *service.ScooterService
	StructureCh    chan *proto.ScooterClient
}

func newHandler(scooterService *service.ScooterService, structure chan *proto.ScooterClient) *handler {
	return &handler{
		scooterService: scooterService,
		StructureCh:    structure,
	}
}

func NewRouter(scooterService *service.ScooterService, structure chan *proto.ScooterClient) *mux.Router {
	router := mux.NewRouter()
	handler := newHandler(scooterService, structure)
	router.HandleFunc(`/scooters`, handler.getAllScooters).Methods("GET")
	router.HandleFunc(`/scooter/{`+scooterIDKey+`}`, handler.getScooterById).Methods("GET")
	router.HandleFunc(`/start-trip/{`+stationIDKey+`}`, handler.showTripPage).Methods("GET")
	router.HandleFunc(`/run`, handler.startScooterTrip).Methods("GET")
	router.HandleFunc(`/choose-station`, handler.chooseStation).Methods("POST")
	router.HandleFunc(`/choose-scooter`, handler.chooseScooter).Methods("POST")
	return router
}

func (h *handler) getAllScooters(w http.ResponseWriter, r *http.Request) {
	scooters, err := h.scooterService.GetAllScooters(context.Background(), &proto.Request{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	json.NewEncoder(w).Encode(scooters)
}

func (h *handler) getScooterById(w http.ResponseWriter, r *http.Request) {
	scooterID, err := strconv.Atoi(mux.Vars(r)[scooterIDKey])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scooter, err := h.scooterService.GetScooterById(context.Background(), &proto.ScooterID{Id: uint64(scooterID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	json.NewEncoder(w).Encode(scooter)
}

func (h *handler) startScooterTrip(w http.ResponseWriter, r *http.Request) {
	//userID is a temporary value
	//userID := 3

	scooterStatus, err := h.scooterService.GetScooterStatus(context.Background(), &proto.ScooterID{Id: uint64(chosenScooterID)})
	if err != nil {
		fmt.Println(err)
	}
	station, err := h.scooterService.GetStationById(context.Background(), &proto.StationID{Id: uint64(chosenStationID)})
	if err != nil {
		fmt.Println(err)
	}

	scooterForClient := proto.ScooterClient{Id: uint64(chosenScooterID), Latitude: scooterStatus.Latitude,
		Longitude: scooterStatus.Longitude, BatteryRemain: scooterStatus.BatteryRemain,
		DestLatitude: station.Latitude, DestLongitude: station.Longitude}

	fmt.Printf("ScooterForClient: %v\n", &scooterForClient)

	h.StructureCh <- &scooterForClient
	fmt.Println("Data has been sent")

	//statusStart, err := h.scooterService.CreateScooterStatusInRent(context.Background(),
	//	&proto.ScooterID{Id: uint64(chosenScooterID)})
	//if err != nil {
	//	fmt.Println(err)
	//}

	//err = h.scooterService.InitAndRun(context.Background(), &proto.ScooterID{Id: uint64(chosenScooterID)},
	//	&proto.StationID{Id: uint64(chosenStationID)})
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//fmt.Println("StatusEnd creating...")
	//statusEnd, err := h.scooterService.CreateScooterStatusInRent(context.Background(),
	//	&proto.ScooterID{Id: uint64(chosenScooterID)})
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Println("StatusEnd created...")
	//tripInfo := &proto.TripInfo{ScooterID: uint64(chosenScooterID), UserID: uint64(userID),
	//	StatusStartID: statusStart.Id,
	//	StatusEndID:   statusEnd.Id}
	//tripOrder, err := h.scooterService.Order.CreateOrder(context.Background(), tripInfo)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(tripOrder)
}

func (h *handler) showTripPage(w http.ResponseWriter, r *http.Request) {
	stationID, err := strconv.Atoi(mux.Vars(r)[stationIDKey])

	scooterList, err := h.scooterService.GetAllScootersByStationID(context.Background(),
		&proto.StationID{Id: uint64(stationID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stationList, err := h.scooterService.GetAllStations(context.Background(), &proto.Request{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//tmpl, err := template.ParseFiles("../scooter_server/templates/scooter-run.html")
	tmpl, err := template.ParseFiles(config.MONO_TEMPLATES_PATH + "scooter-run.html")
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, combineForTemplate{scooterList, stationList})
	if err != nil {
		fmt.Println(err)
	}
}

func (h *handler) chooseScooter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	chosenScooterID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Println(chosenScooterID)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) chooseStation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	chosenStationID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Println(chosenStationID)
	w.WriteHeader(http.StatusOK)
}
