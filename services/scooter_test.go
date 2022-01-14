package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type ScooterMock struct {
	ScooterService *ScooterService
	RepoScooter    *mock.MockScooterRepo
}

type scooterTestCase struct {
	name string
	test func(t *testing.T, mock *ScooterMock)
}

func runScooterTestCases(t *testing.T, testCases []scooterTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := NewScooterMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func NewScooterMock(ctrl *gomock.Controller) *ScooterMock {
	repoScooter := mock.NewMockScooterRepo(ctrl)

	scooterService := NewScooterService(repoScooter)

	return &ScooterMock{
		ScooterService: scooterService,
		RepoScooter:    repoScooter,
	}
}

func TestScooterService_GetScooterById(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().GetScooterById(1).Return(models.ScooterDTO{}, nil).Times(1)

				_, err := mock.ScooterService.GetScooterById(1)
				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().GetScooterById(-1).Return(models.ScooterDTO{}, expectedError).Times(1)

				_, err := mock.ScooterService.GetScooterById(-1)
				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestScooterService_CreateScooterStatusInRent(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().CreateScooterStatusInRent(1).Return(models.ScooterStatusInRent{}, nil).Times(1)

				_, err := mock.RepoScooter.CreateScooterStatusInRent(1)

				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().CreateScooterStatusInRent(1).Return(models.ScooterStatusInRent{}, expectedError).Times(1)

				_, err := mock.ScooterService.CreateScooterStatusInRent(1)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestScooterService_GetAllScooters(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().GetAllScooters().Return(&models.ScooterListDTO{}, nil).Times(1)

				_, err := mock.RepoScooter.GetAllScooters()

				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().GetAllScooters().Return(&models.ScooterListDTO{}, expectedError).Times(1)

				_, err := mock.ScooterService.GetAllScooters()

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestScooterService_GetAllScootersByStationID(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().GetAllScootersByStationID(1).Return(&models.ScooterListDTO{}, nil).Times(1)

				_, err := mock.RepoScooter.GetAllScootersByStationID(1)

				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().GetAllScootersByStationID(1).Return(&models.ScooterListDTO{}, expectedError).Times(1)

				_, err := mock.ScooterService.GetAllScootersByStationID(1)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestScooterService_GetScooterStatus(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().GetScooterStatus(1).Return(models.ScooterStatus{}, nil).Times(1)

				_, err := mock.RepoScooter.GetScooterStatus(1)

				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().GetScooterStatus(1).Return(models.ScooterStatus{}, expectedError).Times(1)

				_, err := mock.ScooterService.GetScooterStatus(1)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}

func TestScooterService_SendCurrentStatus(t *testing.T) {
	runScooterTestCases(t, []scooterTestCase{
		{
			name: "Correct",
			test: func(t *testing.T, mock *ScooterMock) {
				mock.RepoScooter.EXPECT().SendCurrentStatus(1, 1, 40.5, 50.5, 70.0).Return(
					nil).Times(1)

				err := mock.RepoScooter.SendCurrentStatus(1, 1, 40.5, 50.5, 70.0)

				assert.Equal(t, nil, err)
			},
		}, {
			name: "Incorrect",
			test: func(t *testing.T, mock *ScooterMock) {
				expectedError := errors.New("expectedError")
				mock.RepoScooter.EXPECT().SendCurrentStatus(1, 1, 40.5, 50.5, 70.0).Return(expectedError).Times(1)

				err := mock.ScooterService.SendCurrentStatus(1, 1, 40.5, 50.5, 70.0)

				assert.Error(t, err)
				assert.Equal(t, expectedError, err)
			},
		},
	})
}
