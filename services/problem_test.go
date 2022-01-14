package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type problemUseCasesMock struct {
	repoProblem  *mock.MockProblemRepo
	repoSolution *mock.MockSolutionRepo
	problemUC    *ProblemService
}

type problemTestCase struct {
	name string
	test func(t *testing.T, mock *problemUseCasesMock)
}

func runProblemTestCases(t *testing.T, testCases []problemTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := newProblemUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func newProblemUseCasesMock(ctrl *gomock.Controller) *problemUseCasesMock {
	repoProblem := mock.NewMockProblemRepo(ctrl)
	repoSolution := mock.NewMockSolutionRepo(ctrl)
	problemUC := NewProblemService(repoProblem, repoSolution)

	return &problemUseCasesMock{
		repoProblem:  repoProblem,
		repoSolution: repoSolution,
		problemUC:    problemUC,
	}
}

func Test_Problem_GetProblemByID(t *testing.T) {
	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetProblemByID(1).
					Return(models.Problem{ID: 1}, nil).Times(1)

				result, err := mock.problemUC.GetProblemByID(1)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, result.ID)
			},
		},
	})
}

func Test_Problem_AddNewProblem(t *testing.T) {
	modelToReturn := &models.Problem{ID: 1, Description: "Test message"}

	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				mock.repoProblem.EXPECT().AddNewProblem(modelToReturn).
					Return(nil).Times(1)

				err := mock.problemUC.AddNewProblem(modelToReturn)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, modelToReturn.ID)
				assert.Contains(t, modelToReturn.Description, "Test")
			},
		},
	})
}

func Test_Problem_MarkProblemAsSolved(t *testing.T) {
	modelToTest := &models.Problem{ID: 1, IsSolved: false}
	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().MarkProblemAsSolved(modelToTest).
					Return(models.Problem{ID: 1, IsSolved: true}, nil).Times(1)

				result, err := mock.problemUC.MarkProblemAsSolved(modelToTest)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, result.ID)
				assert.True(t, result.IsSolved)
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				someError := errors.New("mark solved problem")

				mock.repoProblem.EXPECT().MarkProblemAsSolved(modelToTest).
					Return(models.Problem{ID: 1, IsSolved: false}, someError).Times(1)

				result, err := mock.problemUC.MarkProblemAsSolved(modelToTest)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
				assert.False(t, result.IsSolved)
			},
		},
	})
}

func Test_Problem_GetProblemTypeByID(t *testing.T) {
	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetProblemTypeByID(1).
					Return(models.ProblemType{ID: 1}, nil).Times(1)

				result, err := mock.problemUC.GetProblemTypeByID(1)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, result.ID)
			},
		},
	})
}

func Test_Problem_GetAllProblemTypes(t *testing.T) {
	modelsToReturn := []models.ProblemType{{ID: 1, Name: "Test1"}, {ID: 2, Name: "Test2"}}

	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetAllProblemTypes().
					Return(modelsToReturn, nil).Times(1)

				result, err := mock.problemUC.GetAllProblemTypes()
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result))
				assert.Contains(t, result[0].Name, "Test")
			},
		},
	})
}

func Test_Problem_GetProblemsByUserID(t *testing.T) {
	modelsToReturn := &models.ProblemList{
		Problems: []models.Problem{
			{ID: 1, User: models.User{ID: 1}, Description: "Test1"},
			{ID: 2, User: models.User{ID: 1}, Description: "Test2"},
		},
	}

	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetProblemsByUserID(1).
					Return(modelsToReturn, nil).Times(1)

				result, err := mock.problemUC.GetProblemsByUserID(1)
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result.Problems))
				assert.Contains(t, result.Problems[0].Description, "Test")
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				someError := errors.New("get problem list by user")

				mock.repoProblem.EXPECT().GetProblemsByUserID(2).
					Return(&models.ProblemList{}, someError).Times(1)

				result, err := mock.problemUC.GetProblemsByUserID(2)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
				assert.Empty(t, result.Problems)
			},
		},
	})
}

func Test_Problem_GetProblemsByTypeID(t *testing.T) {
	modelsToReturn := &models.ProblemList{
		Problems: []models.Problem{
			{ID: 1, Type: models.ProblemType{ID: 1}, Description: "Test1"},
			{ID: 2, Type: models.ProblemType{ID: 1}, Description: "Test2"},
		},
	}

	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetProblemsByTypeID(1).
					Return(modelsToReturn, nil).Times(1)

				result, err := mock.problemUC.GetProblemsByTypeID(1)
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result.Problems))
				assert.Contains(t, result.Problems[0].Description, "Test")
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				someError := errors.New("get problem list by type")

				mock.repoProblem.EXPECT().GetProblemsByTypeID(2).
					Return(&models.ProblemList{}, someError).Times(1)

				result, err := mock.problemUC.GetProblemsByTypeID(2)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
				assert.Empty(t, result.Problems)
			},
		},
	})
}

func Test_Problem_AddProblemSolution(t *testing.T) {
	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoSolution.EXPECT().AddProblemSolution(1, &models.Solution{}).
					Return(nil).Times(1)
				mock.repoProblem.EXPECT().GetProblemByID(1).
					Return(models.Problem{IsSolved: false}, nil).Times(1)
				mock.repoProblem.EXPECT().MarkProblemAsSolved(&models.Problem{IsSolved: false}).
					Return(models.Problem{IsSolved: true}, nil).Times(1)

				err := mock.problemUC.AddProblemSolution(1, &models.Solution{})
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				someError := errors.New("add problem solution")

				mock.repoSolution.EXPECT().AddProblemSolution(1, &models.Solution{}).
					Return(someError).Times(1)

				err := mock.problemUC.AddProblemSolution(1, &models.Solution{})
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
	})
}
