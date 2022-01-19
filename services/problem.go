package services

import (
	"Dp218GO/models"
	"context"
	"google.golang.org/grpc"
	proto2 "problem.micro/proto"
	"time"
)

// ProblemService - structure for implementing user problem service
type ProblemService struct {
	microservice proto2.ProblemServiceClient
	userService  *UserService
}

func (problserv *ProblemService) unmarshallProblem(problemGRPC *proto2.Problem) models.Problem {
	problem := models.Problem{
		ID:           int(problemGRPC.Id),
		DateReported: time.Unix(problemGRPC.ReportedAt.Seconds, 0),
		Description:  problemGRPC.Description,
		IsSolved:     problemGRPC.IsSolved,
	}
	problserv.AddProblemComplexFields(&problem, int(problemGRPC.Type.Id), int(problemGRPC.UserId))
	return problem
}

func (problserv *ProblemService) unmarshallProblemType(problemTypeGRPC *proto2.ProblemType) models.ProblemType {
	return models.ProblemType{
		ID:   int(problemTypeGRPC.Id),
		Name: problemTypeGRPC.Name,
	}
}

func (problserv *ProblemService) unmarshallSolution(solutionGRPC *proto2.Solution) models.Solution {
	solution := models.Solution{
		Problem:     problserv.unmarshallProblem(solutionGRPC.Problem),
		DateSolved:  time.Unix(solutionGRPC.SolvedAt.Seconds, 0),
		Description: solutionGRPC.Description,
	}

	return solution
}

func (problserv *ProblemService) marshallProblem(problem *models.Problem) *proto2.Problem {
	return &proto2.Problem{
		Id:          int64(problem.ID),
		UserId:      int64(problem.User.ID),
		Description: problem.Description,
		Type:        &proto2.ProblemType{Id: int32(problem.Type.ID), Name: problem.Type.Name},
		IsSolved:    problem.IsSolved,
		ReportedAt:  &proto2.DateTime{Seconds: problem.DateReported.Unix()},
	}
}

func (problserv *ProblemService) marshallSolution(solution *models.Solution) *proto2.Solution {
	return &proto2.Solution{
		Problem:     problserv.marshallProblem(&solution.Problem),
		Description: solution.Description,
		SolvedAt:    &proto2.DateTime{Seconds: solution.DateSolved.Unix()},
	}
}

// NewProblemService - initialization of ProblemService
func NewProblemService(grpcConn grpc.ClientConnInterface, userServ *UserService) *ProblemService {
	return &ProblemService{
		microservice: proto2.NewProblemServiceClient(grpcConn),
		userService:  userServ,
	}
}

// AddNewProblem - add new user problem record
func (problserv *ProblemService) AddNewProblem(problem *models.Problem) error {
	problemType := &proto2.ProblemType{
		Id: int32(problem.Type.ID),
	}
	problemToAdd := &proto2.Problem{
		UserId:      int64(problem.User.ID),
		Description: problem.Description,
		Type:        problemType,
		IsSolved:    problem.IsSolved,
	}
	_, err := problserv.microservice.AddNewProblem(context.Background(), problemToAdd)
	return err
}

// GetProblemByID - get problem information by its ID
func (problserv *ProblemService) GetProblemByID(problemID int) (models.Problem, error) {
	request := &proto2.ProblemRequest{Id: int64(problemID)}
	response, err := problserv.microservice.GetProblemByID(context.Background(), request)

	return problserv.unmarshallProblem(response.Problem), err
}

// MarkProblemAsSolved - update problem record to make problem solved
func (problserv *ProblemService) MarkProblemAsSolved(problem *models.Problem) (models.Problem, error) {
	problem.IsSolved = true
	response, err := problserv.microservice.UpdateProblem(context.Background(), problserv.marshallProblem(problem))

	return problserv.unmarshallProblem(response.Problem), err
}

// GetProblemTypeByID - get problem type record by its ID
func (problserv *ProblemService) GetProblemTypeByID(typeID int) (models.ProblemType, error) {
	request := &proto2.ProblemRequest{TypeId: int32(typeID)}
	response, err := problserv.microservice.GetProblemTypeByID(context.Background(), request)
	return problserv.unmarshallProblemType(response.ProblemType), err
}

func (problserv *ProblemService) GetAllProblemTypes() ([]models.ProblemType, error) {
	request := &proto2.ProblemRequest{}
	response, err := problserv.microservice.GetAllProblemTypes(context.Background(), request)
	var resultingList []models.ProblemType
	if err != nil {
		return resultingList, err
	}
	for _, val := range response.ProblemTypes {
		resultingList = append(resultingList, problserv.unmarshallProblemType(val))
	}
	return resultingList, err
}

// GetProblemsByUserID - get problem list for given user (by user ID)
func (problserv *ProblemService) GetProblemsByUserID(userID int) (*models.ProblemList, error) {
	request := &proto2.ProblemRequest{UserId: int64(userID)}
	response, err := problserv.microservice.GetProblemsByUserID(context.Background(), request)
	resultingList := &models.ProblemList{}
	if err != nil {
		return resultingList, err
	}
	for _, val := range response.Problems {
		resultingList.Problems = append(resultingList.Problems, problserv.unmarshallProblem(val))
	}
	return resultingList, err
}

// GetProblemsByTypeID - get problem list by given problem type ID
func (problserv *ProblemService) GetProblemsByTypeID(typeID int) (*models.ProblemList, error) {
	request := &proto2.ProblemRequest{TypeId: int32(typeID)}
	response, err := problserv.microservice.GetProblemsByTypeID(context.Background(), request)
	resultingList := &models.ProblemList{}
	if err != nil {
		return resultingList, err
	}
	for _, val := range response.Problems {
		resultingList.Problems = append(resultingList.Problems, problserv.unmarshallProblem(val))
	}
	return resultingList, err
}

// GetProblemsByBeingSolved - get problem list by is_solved field value
func (problserv *ProblemService) GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error) {
	request := &proto2.ProblemRequest{IsSolved: solved}
	response, err := problserv.microservice.GetProblemsBySolved(context.Background(), request)
	resultingList := &models.ProblemList{}
	if err != nil {
		return resultingList, err
	}
	for _, val := range response.Problems {
		resultingList.Problems = append(resultingList.Problems, problserv.unmarshallProblem(val))
	}
	return resultingList, err
}

// GetProblemsByTimePeriod - get problem list from time start to time end
func (problserv *ProblemService) GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error) {
	request := &proto2.ProblemRequest{
		StartTime: &proto2.DateTime{Seconds: start.Unix()},
		EndTime:   &proto2.DateTime{Seconds: end.Unix()},
	}
	response, err := problserv.microservice.GetProblemsByTimePeriod(context.Background(), request)
	resultingList := &models.ProblemList{}
	if err != nil {
		return resultingList, err
	}
	for _, val := range response.Problems {
		resultingList.Problems = append(resultingList.Problems, problserv.unmarshallProblem(val))
	}
	return resultingList, err
}

// AddProblemComplexFields - fulfill problem model with problem type, scooter, user (by their IDs)
func (problserv *ProblemService) AddProblemComplexFields(problem *models.Problem, typeID, userID int) {
	if typeID != 0 {
		problemType, eType := problserv.GetProblemTypeByID(typeID)
		if eType == nil {
			problem.Type = problemType
		}
	}
	if userID != 0 {
		user, eUser := problserv.userService.GetUserByID(userID)
		if eUser == nil {
			problem.User = user
		}
	}
}

// AddProblemSolution - make solution record for given problem (by ID)
func (problserv *ProblemService) AddProblemSolution(problemID int, solution *models.Solution) error {
	request := &proto2.ProblemSolution{
		Problem:  &proto2.Problem{Id: int64(problemID)},
		Solution: problserv.marshallSolution(solution),
	}
	_, err := problserv.microservice.AddProblemSolution(context.Background(), request)
	return err
}

// GetSolutionByProblem - get solution for given problem
func (problserv *ProblemService) GetSolutionByProblem(problem models.Problem) (models.Solution, error) {
	response, err := problserv.microservice.GetSolutionByProblem(context.Background(), problserv.marshallProblem(&problem))

	return problserv.unmarshallSolution(response.Solution), err
}
