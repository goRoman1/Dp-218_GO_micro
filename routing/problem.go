package routing

import (
	"Dp218GO/models"
	"Dp218GO/services"
	"Dp218GO/utils"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

var problemService *services.ProblemService
var problemIDKey = "problemID"

var keyProblemRoutes = []Route{
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

type problemsForTemplate struct {
	ProblemList *models.ProblemList
}

// DistinctProblemUsers - to fill users in templates filter
func (pu *problemsForTemplate) DistinctProblemUsers() map[int]models.User {
	var result = make(map[int]models.User)
	for _, v := range pu.ProblemList.Problems {
		result[v.User.ID] = v.User
	}
	return result
}

// DistinctProblemTypes - to fill types of problems in templates filter
func (pu *problemsForTemplate) DistinctProblemTypes() map[int]models.ProblemType {
	var result = make(map[int]models.ProblemType)
	for _, v := range pu.ProblemList.Problems {
		result[v.Type.ID] = v.Type
	}
	return result
}

// AddProblemHandler - add endpoints for user problems & solutions to http router
func AddProblemHandler(router *mux.Router, problserv *services.ProblemService) {
	problemService = problserv
	problemRouter := router.NewRoute().Subrouter()
	problemRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keyProblemRoutes {
		problemRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		problemRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllProblems(w http.ResponseWriter, r *http.Request) {

	var problems *models.ProblemList
	var err error
	var userID, typeID, dateFrom, dateTo, isSolvedFilter interface{}
	format := GetFormatFromRequest(r)

	userID, err = GetParameterFromRequest(r, "UserID", utils.ConvertStringToInt())
	if err == nil {
		problems, err = problemService.GetProblemsByUserID(userID.(int))
		if err != nil {
			ServerErrorRender(format, w)
			return
		}
	}

	if err != nil {
		typeID, err = GetParameterFromRequest(r, "TypeID", utils.ConvertStringToInt())
		if err == nil {
			problems, err = problemService.GetProblemsByTypeID(typeID.(int))
			if err != nil {
				ServerErrorRender(format, w)
				return
			}
		}
	}

	if err != nil {
		dateFrom, err = GetParameterFromRequest(r, "DateFrom", utils.ConvertStringToTime())
		if err == nil {
			dateTo, err = GetParameterFromRequest(r, "DateTo", utils.ConvertStringToTime())
			if err == nil {
				problems, err = problemService.GetProblemsByTimePeriod(dateFrom.(time.Time), dateTo.(time.Time))
				if err != nil {
					ServerErrorRender(format, w)
					return
				}
			}
		}
	}

	if err != nil {
		isSolvedFilter, err = GetParameterFromRequest(r, "SolvedFilter", utils.ConvertStringToBool())
		if err == nil {
			problems, err = problemService.GetProblemsByBeingSolved(isSolvedFilter.(bool))
			if err != nil {
				ServerErrorRender(format, w)
				return
			}
		}
	}

	if err != nil {
		problems, err = problemService.GetProblemsByTimePeriod(time.Unix(0, 0), time.Now())
		if err != nil {
			ServerErrorRender(format, w)
			return
		}
	}

	EncodeAnswer(format, w, &problemsForTemplate{problems}, HTMLPath+"problems.html")
}

func newProblem(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(errors.New("not authorized")))
		return
	}

	problemTypes, err := problemService.GetAllProblemTypes()
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	problem := &models.Problem{User: *user}
	problemWithAllTypes := struct {
		Problem *models.Problem
		Types   []models.ProblemType
	}{problem, problemTypes}

	EncodeAnswer(FormatHTML, w, problemWithAllTypes, HTMLPath+"problem-add.html")
}

func getProblemInfo(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	problemID, err := strconv.Atoi(mux.Vars(r)[problemIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	problem, err := problemService.GetProblemByID(problemID)
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, problem, HTMLPath+"problem.html")
}

func addProblem(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	problemData := models.Problem{}
	DecodeRequest(format, w, r, &problemData, decodeProblemAddRequest)
	err := problemService.AddNewProblem(&problemData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	if format == FormatHTML {
		getAllProblems(w, r)
		return
	}
	EncodeAnswer(FormatJSON, w, problemData)
}

func decodeProblemAddRequest(r *http.Request, data interface{}) error {

	var err error

	problemData := data.(*models.Problem)

	description, _ := GetParameterFromRequest(r, "Description", utils.ConvertStringToString())
	userID, err := GetParameterFromRequest(r, "UserID", utils.ConvertStringToInt())
	if err != nil {
		return err
	}
	typeID, err := GetParameterFromRequest(r, "TypeID", utils.ConvertStringToInt())
	if err != nil {
		return err
	}

	problemData.Description = description.(string)
	problemData.IsSolved = false
	problemService.AddProblemComplexFields(problemData, typeID.(int), userID.(int))

	return err
}

func getProblemSolution(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	problemID, err := strconv.Atoi(mux.Vars(r)[problemIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	problem, err := problemService.GetProblemByID(problemID)
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	solution, err := problemService.GetSolutionByProblem(problem)
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, solution, HTMLPath+"solution.html")
}

func addProblemSolution(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	problemID, err := strconv.Atoi(mux.Vars(r)[problemIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	solutionData := models.Solution{}
	solutionData.Problem = models.Problem{ID: problemID}
	DecodeRequest(format, w, r, &solutionData, decodeSolutionAddRequest)
	err = problemService.AddProblemSolution(solutionData.Problem.ID, &solutionData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	getProblemInfo(w, r)
}

func decodeSolutionAddRequest(r *http.Request, data interface{}) error {
	var err error

	solutionData := data.(*models.Solution)

	description, _ := GetParameterFromRequest(r, "Description", utils.ConvertStringToString())
	if err != nil {
		return err
	}
	solutionData.Description = description.(string)

	problemID, err := GetParameterFromRequest(r, "ProblemID", utils.ConvertStringToInt())
	if err != nil {
		return err
	}
	problem, err := problemService.GetProblemByID(problemID.(int))
	if err != nil {
		return err
	}
	solutionData.Problem = problem

	return problemService.AddProblemSolution(problemID.(int), solutionData)
}
