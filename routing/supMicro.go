package routing

import (
	"Dp218GO/models"
	"Dp218GO/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var supMicroService *services.SupMicroService

var keySupMicroRoutes = []Route{
	{
		Uri:     `/mStations`,
		Method:  http.MethodGet,
		Handler: getTempData,
	},
	{
		Uri:     `/microAddStation`,
		Method:  http.MethodPost,
		Handler: createMicroStation,
	},
	{
		Uri:     `/microASl`,
		Method:  http.MethodPost,
		Handler: createMicroStationInLocation,
	},
}

// AddSupMicroHandler - add endpoints for user problems & solutions to http router
func AddSupMicroHandler(router *mux.Router, supserv *services.SupMicroService) {
	supMicroService = supserv
	supMicroRouter := router.NewRoute().Subrouter()

	for _, rt := range keySupMicroRoutes {
		supMicroRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		supMicroRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getTempData(w http.ResponseWriter, r *http.Request) {

	var err error
	var station = &models.StationList{}
	var location = &models.LocationList{}
	format := GetFormatFromRequest(r)

	station, err = stationService.GetAllStations()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	location, err = supMicroService.GetLocations()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	list := &models.StationLocation{}

	list.StationList = *station
	list.LocationList = *location

	EncodeAnswer(format, w, list, HTMLPath+"supplierMicro.html")
}

func createMicroStationInLocation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	stationData := models.Station{}
	locationData := models.Location{}
	err := supplierMicroService.CreateStationInLocation(&locationData, &stationData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(FormatJSON, w, stationData)
}

func createMicroStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	stationData := models.Station{}

	name := r.FormValue("stationName")
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")

	fLatitude, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Println(err)
	}
	fLongitude, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Println(err)
	}

	stationData.Name = name
	stationData.IsActive = true
	stationData.Latitude = fLatitude
	stationData.Longitude = fLongitude

	if err := supplierMicroService.AddNewStation(&stationData); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/mStations", http.StatusFound)
}
