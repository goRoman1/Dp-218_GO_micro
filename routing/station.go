package routing

import (
	"Dp218GO/models"
	"fmt"
	"net/http"
	"strconv"

	"Dp218GO/services"
	"github.com/gorilla/mux"
)

var stationService *services.StationService
var stationIDKey = "stationID"

var keyRoutesStation = []Route{
	{
		Uri:     `/stations`,
		Method:  http.MethodGet,
		Handler: getAllStations,
	},
	{
		Uri:     `/station/{` + stationIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getStation,
	},
	{
		Uri:     `/station`,
		Method:  http.MethodPost,
		Handler: createStation,
	},
	{
		Uri:     `/station/{` + stationIDKey + `}`,
		Method:  http.MethodDelete,
		Handler: deleteStation,
	},
	{
		Uri:     `/stations`,
		Method:  http.MethodPost,
		Handler: allStationsOperation,
	},
	{
		Uri:     `/station/{` + stationIDKey + `}`,
		Method:  http.MethodPost,
		Handler: UpdateStation,
	},
}

func AddStationHandler(router *mux.Router, service *services.StationService) {
	stationService = service
	for _, rt := range keyRoutesStation {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	station := &models.Station{}
	DecodeRequest(format, w, r, station, nil)

	if err := stationService.AddStation(station); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, station)
}

func getAllStations(w http.ResponseWriter, r *http.Request) {
	var station = &models.StationList{}
	var err error
	format := GetFormatFromRequest(r)

	station, err = stationService.GetAllStations()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, station, HTMLPath+"station-list.html")
}

func getStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationId, err := strconv.Atoi(mux.Vars(r)[stationIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	station, err := stationService.GetStationById(stationId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, station, HTMLPath+"station-edit.html")
}

func deleteStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationId, err := strconv.Atoi(mux.Vars(r)[stationIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = stationService.DeleteStation(stationId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}

func allStationsOperation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	r.ParseForm()
	if _, ok := r.Form["ActionType"]; !ok {

		return
	}
	actionType := r.FormValue("ActionType")
	switch actionType {
	case "BlockStation":
		stationId, err := strconv.Atoi(r.FormValue("stationID"))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = stationService.ChangeStationBlockStatus(stationId)
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		EncodeError(format, w, ErrorRendererDefault(fmt.Errorf("unknown users operation")))
	}
	getAllStations(w, r)
}

func UpdateStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationId, err := strconv.Atoi(mux.Vars(r)[stationIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	stationData, err := stationService.GetStationById(stationId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	DecodeRequest(format, w, r, &stationData, DecodeStationUpdateRequest)
	stationData, err = stationService.UpdateStation(stationId, stationData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, stationData, HTMLPath+"station-edit.html")
}

func DecodeStationUpdateRequest(r *http.Request, data interface{}) error {
	r.ParseForm()
	stationData := data.(*models.Station)
	if _, ok := r.Form["IsActive"]; ok {
		stationData.IsActive, _ = strconv.ParseBool(r.FormValue("IsActive"))
	}
	if _, ok := r.Form["Name"]; ok {
		stationData.Name = r.FormValue("Name")
	}
	if _, ok := r.Form["Latitude"]; ok {
		stationData.Latitude, _ = strconv.ParseFloat(r.FormValue("Latitude"), 64)
	}
	if _, ok := r.Form["Longitude"]; ok {
		stationData.Longitude, _ = strconv.ParseFloat(r.FormValue("Longitude"), 64)
	}
	return nil
}
