package routing

import (
	"Dp218GO/models"
	"Dp218GO/services"
	"github.com/gorilla/mux"
	"net/http"
)

var supplierMicroService *services.SupplierMicroService
var supplierMicroIDKey = "ID"

var keySupplierMicroRoutes = []Route{
	{
		Uri:     `/problems`,
		Method:  http.MethodGet,
		Handler: getAllProblems,
	},
	{
		Uri:     `/problem/{` + problemIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getProblemInfo,
	},
	{
		Uri:     `/problems`,
		Method:  http.MethodPost,
		Handler: addProblem,
	},
	{
		Uri:     `/problem`,
		Method:  http.MethodGet,
		Handler: newProblem,
	},
	{
		Uri:     `/problem/{` + problemIDKey + `}/solution`,
		Method:  http.MethodPost,
		Handler: addProblemSolution,
	},
	{
		Uri:     `/problem/{` + problemIDKey + `}/solution`,
		Method:  http.MethodGet,
		Handler: getProblemSolution,
	},
}

type supplierMicroForTemplate struct {
	ProblemList *models.ProblemList
}

// DistinctProblemUsers - to fill users in templates filter
func (pu *problemsForTemplate) DistinctSupplierMicroUsers() map[int]models.User {
	var result = make(map[int]models.User)
	for _, v := range pu.ProblemList.Problems {
		result[v.User.ID] = v.User
	}
	return result
}

// DistinctProblemTypes - to fill types of problems in templates filter
func (pu *problemsForTemplate) DistinctSupplierMicroTypes() map[int]models.ProblemType {
	var result = make(map[int]models.ProblemType)
	for _, v := range pu.ProblemList.Problems {
		result[v.Type.ID] = v.Type
	}
	return result
}

// AddProblemHandler - add endpoints for user problems & solutions to http router
func AddSupplierMicroHandler(router *mux.Router, supserv *services.SupplierMicroService) {
	supplierMicroService = supserv
	supplierMicroRouter := router.NewRoute().Subrouter()
	supplierMicroRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keyProblemRoutes {
		supplierMicroRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		supplierMicroRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}
