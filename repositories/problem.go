//go:generate mockgen -source=problem.go -destination=../repositories/mock/mock_problem.go -package=mock
package repositories

import (
	"Dp218GO/models"
	"time"
)

// ProblemRepo - interface for user problem repository
type ProblemRepo interface {
	AddNewProblem(problem *models.Problem) error
	GetProblemByID(problemID int) (models.Problem, error)
	GetProblemTypeByID(typeID int) (models.ProblemType, error)
	GetProblemsByUserID(userID int) (*models.ProblemList, error)
	GetProblemsByTypeID(typeID int) (*models.ProblemList, error)
	GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error)
	GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error)
	AddProblemComplexFields(problem *models.Problem, typeID, userID int)
	MarkProblemAsSolved(problem *models.Problem) (models.Problem, error)
	GetAllProblemTypes() ([]models.ProblemType, error)
}

// SolutionRepo - interface for solution repository
type SolutionRepo interface {
	AddProblemSolution(problemID int, solution *models.Solution) error
	GetSolutionByProblem(problem models.Problem) (models.Solution, error)
}
