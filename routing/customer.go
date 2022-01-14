package routing

import (
	"Dp-218_GO_micro/internal/validation"
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type customerHandler struct {
	custService *services.CustomerService
}

func newCustomerHandler(service *services.CustomerService) *customerHandler {
	return &customerHandler{
		custService: service,
	}
}

//AddCustomerHandler registeres endpoints for customer
func AddCustomerHandler(router *mux.Router, service *services.CustomerService) {

	custHandler := newCustomerHandler(service)

	custRouter := router.PathPrefix("/customer").Subrouter()
	custRouter.Use(FilterAuth(authenticationService), FilterCustomer)

	custRouter.Path("/map").HandlerFunc(custHandler.HomeHandler).Methods(http.MethodGet)
	custRouter.Path("/station").HandlerFunc(custHandler.StationListHandler).Methods(http.MethodGet)
	custRouter.Path("/station/nearest").
		HandlerFunc(custHandler.StationNearestHandler).Queries("x", "{x}", "y", "{y}").Methods(http.MethodGet)
	custRouter.Path("/station/{id:[0-9]+}").HandlerFunc(custHandler.StationInfoHandler).Methods(http.MethodGet)

}

// HomeHandler is handler for rendering home page of customer
func (h *customerHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	format := GetFormatFromRequest(r)

	// no need if wrapped with filteruser
	if user == nil {
		EncodeError(format, w, ErrorRendererDefault(errors.New("not authenticated")))
		return
	}

	EncodeAnswer(format, w, user, HTMLPath+"customer-map.html")
}

// StationListHandler is handler that users customer service
// shows list of available stations on map
// returns json station list in response shows error if failed
func (h *customerHandler) StationListHandler(w http.ResponseWriter, r *http.Request) {
	// TODO show only not blocked
	sts, err := h.custService.ListStations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sts.Station)
}

//StationNearestHandler is handler that user customer service
// takes user location and returns nearest
// station in json format shows error if failed
func (h *customerHandler) StationNearestHandler(w http.ResponseWriter, r *http.Request) {

	valReq := validation.LocationRequest{
		Latitude:  r.FormValue("x"),
		Longitude: r.FormValue("y"),
	}

	if err := valReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	x, err := strconv.ParseFloat(valReq.Latitude, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	y, err := strconv.ParseFloat(valReq.Longitude, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nearest, err := h.custService.ShowNearestStation(x, y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]*models.Station{nearest})
}

// StationInfoHandler is handler that shows general station info of station
// by station id received in reguest url var
func (h *customerHandler) StationInfoHandler(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	station, err := h.custService.ShowStation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)
}

// FilterCustomer is middleware that restricts access to customer page
// checks if user role is customer or admin
// shows error if not allowed
func FilterCustomer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user == nil || !(user.Role.IsUser || user.Role.IsAdmin) {
			http.Error(w, "only customers allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
