package routing

import (
	"Dp218GO/models"
	"Dp218GO/services"
	"github.com/gorilla/mux"
	"net/http"
)

var supplierMicroService *services.SupplierMicroService

var keySupplierMicroRoutes = []Route{
	{
		Uri:     `/render`,
		Method:  http.MethodGet,
		Handler: getTemplateData,
	}, {
		Uri:     `/locations`,
		Method:  http.MethodGet,
		Handler: getLocations,
	},
	{
		Uri:     `/addStation`,
		Method:  http.MethodPost,
		Handler: addStation,
	},
	{
		Uri:     `/ASl`,
		Method:  http.MethodPost,
		Handler: addStationInLocation,
	},
}

// AddProblemHandler - add endpoints for user problems & solutions to http router
func AddSupplierMicroHandler(router *mux.Router, supserv *services.SupplierMicroService) {
	supplierMicroService = supserv
	supplierMicroRouter := router.NewRoute().Subrouter()
	//	supplierMicroRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keySupplierMicroRoutes {
		supplierMicroRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		supplierMicroRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getTemplateData(w http.ResponseWriter, r *http.Request) {
	var err error
	var locationData *models.LocationList
	var stationsData models.StationList
	format := GetFormatFromRequest(r)
	locationData, err = supplierMicroService.GetAllLocations()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	stationsData, err = supplierMicroService.GetAllLStations()
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	problemWithAllTypes := struct {
		locationList *models.LocationList
		stationList  models.StationList
	}{locationData, stationsData}

	EncodeAnswer(FormatHTML, w, problemWithAllTypes, HTMLPath+"supplierMicro.html")
}

func addStationInLocation(w http.ResponseWriter, r *http.Request) {
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

func addStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationData := models.Station{}
	DecodeRequest(format, w, r, &stationData, decodeProblemAddRequest)
	err := supplierMicroService.AddNewStation(&stationData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	if format == FormatHTML {
		getTemplateData(w, r)
		return
	}
	EncodeAnswer(FormatJSON, w, stationData)
}

func getLocations(w http.ResponseWriter, r *http.Request) {
	//	format := GetFormatFromRequest(r)
	var locationData *models.LocationList

	locationData, err := supplierMicroService.GetAllLocations()
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatHTML, w, locationData, HTMLPath+"supplierMicro.html")
}
