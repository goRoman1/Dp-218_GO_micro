package routing

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/services"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var scooterInitService *services.ScooterInitService
var scooterInitRoutes = []Route{
	{
		Uri:     `/init`,
		Method:  http.MethodGet,
		Handler: getAllocationData,
	},
	{
		Uri:     `/transfer`,
		Method:  http.MethodPost,
		Handler: addStatusesToScooters,
	},
}

// AddScooterInitHandler - add endpoints for working with scooterInit to http router
func AddScooterInitHandler(router *mux.Router, service *services.ScooterInitService) {
	scooterInitService = service
	scooterInitRouter := router.NewRoute().Subrouter()
	scooterInitRouter.Use(FilterAuth(authenticationService))

	for _, rt := range scooterInitRoutes {
		scooterInitRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		scooterInitRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

// getAllocationData - render data on page
func getAllocationData(w http.ResponseWriter, r *http.Request) {
	var dataAllocation = &models.ScootersStationsAllocation{}
	var err error
	format := GetFormatFromRequest(r)
	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	dataAllocation = scooterInitService.ConvertForTemplateStruct()

	EncodeAnswer(format, w, dataAllocation, HTMLPath+"scooter-init.html")
}

// addStatusesToScooters - add statuses to scooter statuses
func addStatusesToScooters(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	scooterIds := r.Form["new_data"]
	stationId := r.Form["station_data"]
	intStationId, err := strconv.Atoi(stationId[0])
	if err != nil {
		log.Println(err)
	}

	var intScooterIds []int
	for _, i := range scooterIds {
		intId, err := strconv.Atoi(i)
		if err != nil {
			log.Println(err)
		}
		intScooterIds = append(intScooterIds, intId)
	}

	stationData, err := stationService.GetStationById(intStationId)

	err = scooterInitService.AddStatusesToScooters(intScooterIds, stationData)
	if err != nil {
		return
	}

	http.Redirect(w, r, "http://localhost:8080/init", http.StatusFound)
}
