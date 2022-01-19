package services

import (
	"Dp218GO/models"
	mock "Dp218GO/repositories/mock"
	"errors"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type userUseCasesMock struct {
	repoUser *mock.MockUserRepo
	repoRole *mock.MockRoleRepo
	userUC   *UserService
}

type userTestCase struct {
	name string
	test func(t *testing.T, mock *userUseCasesMock)
}

func runUserTestCases(t *testing.T, testCases []userTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := newUserUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func newUserUseCasesMock(ctrl *gomock.Controller) *userUseCasesMock {
	repoUser := mock.NewMockUserRepo(ctrl)
	repoRole := mock.NewMockRoleRepo(ctrl)
	userUC := NewUserService(repoUser, repoRole)

	return &userUseCasesMock{
		repoUser: repoUser,
		repoRole: repoRole,
		userUC:   userUC,
	}
}

func Test_User_ChangeUsersBlockStatus(t *testing.T) {
	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().GetUserByID(1).
					Return(models.User{IsBlocked: false}, nil).Times(1)

				mock.repoUser.EXPECT().UpdateUser(1, models.User{IsBlocked: true}).
					Return(models.User{}, nil).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(1)
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "incorrect get by ID",
			test: func(t *testing.T, mock *userUseCasesMock) {

				var someError = errors.New("error get by ID")

				mock.repoUser.EXPECT().GetUserByID(2).
					Return(models.User{IsBlocked: false}, someError).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(2)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
		{
			name: "incorrect update user",
			test: func(t *testing.T, mock *userUseCasesMock) {

				var someError = errors.New("error update user")

				mock.repoUser.EXPECT().GetUserByID(3).
					Return(models.User{IsBlocked: false}, nil).Times(1)

				mock.repoUser.EXPECT().UpdateUser(3, models.User{IsBlocked: true}).
					Return(models.User{}, someError).Times(1)

				err := mock.userUC.ChangeUsersBlockStatus(3)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
	})
}

func Test_User_GetUserByID(t *testing.T) {
	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().GetUserByID(1).
					Return(models.User{}, nil).Times(1)

				_, err := mock.userUC.GetUserByID(1)
				assert.Equal(t, nil, err)
			},
		},
	})
}

func Test_User_UpdateUser(t *testing.T) {
	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().UpdateUser(1, models.User{UserName: "Test"}).
					Return(models.User{UserName: "Test"}, nil).Times(1)

				result, err := mock.userUC.UpdateUser(1, models.User{UserName: "Test"})
				assert.Equal(t, nil, err)
				assert.Equal(t, "Test", result.UserName)
			},
		},
	})
}

func Test_User_DeleteUser(t *testing.T) {
	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().DeleteUser(1).
					Return(nil).Times(1)

				err := mock.userUC.DeleteUser(1)
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *userUseCasesMock) {

				someError := errors.New("deletion error")

				mock.repoUser.EXPECT().DeleteUser(1).
					Return(someError).Times(1)

				err := mock.userUC.DeleteUser(1)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
	})
}

func Test_User_AddUser(t *testing.T) {
	modelToReturn := &models.User{ID: 1, UserName: "Test"}

	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().AddUser(modelToReturn).
					Return(nil).Times(1)

				err := mock.userUC.AddUser(modelToReturn)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, modelToReturn.ID)
				assert.Equal(t, "Test", modelToReturn.UserName)
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *userUseCasesMock) {

				someError := errors.New("addition error")

				mock.repoUser.EXPECT().AddUser(modelToReturn).
					Return(someError).Times(1)

				err := mock.userUC.AddUser(modelToReturn)
				assert.Error(t, err)
				assert.Equal(t, someError, err)
			},
		},
	})
}

func Test_User_GetAllUsers(t *testing.T) {
	modelsToReturn := &models.UserList{Users: []models.User{{ID: 1, UserName: "Test1"}, {ID: 2, UserName: "Test2"}}}

	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().GetAllUsers().
					Return(modelsToReturn, nil).Times(1)

				result, err := mock.userUC.GetAllUsers()
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result.Users))
				assert.Contains(t, result.Users[0].UserName, "Test")
			},
		},
	})
}

func Test_User_FindUsersByLoginNameSurname(t *testing.T) {
	modelsToReturn := &models.UserList{Users: []models.User{{ID: 1, UserName: "Test1"}, {ID: 2, UserName: "Test2"}}}

	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoUser.EXPECT().FindUsersByLoginNameSurname("Test").
					Return(modelsToReturn, nil).Times(1)

				result, err := mock.userUC.FindUsersByLoginNameSurname("Test")
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result.Users))
				assert.Contains(t, result.Users[0].UserName, "Test")
			},
		},
	})
}

func Test_User_GetAllRoles(t *testing.T) {
	rolesToReturn := &models.RoleList{Roles: []models.Role{{ID: 1, Name: "Test1"}, {ID: 2, Name: "Test2"}}}

	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoRole.EXPECT().GetAllRoles().
					Return(rolesToReturn, nil).Times(1)

				result, err := mock.userUC.GetAllRoles()
				assert.Equal(t, nil, err)
				assert.Equal(t, 2, len(result.Roles))
				assert.Contains(t, result.Roles[0].Name, "Test")
			},
		},
	})
}

func Test_User_GetRoleByID(t *testing.T) {
	runUserTestCases(t, []userTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *userUseCasesMock) {

				mock.repoRole.EXPECT().GetRoleByID(1).
					Return(models.Role{}, nil).Times(1)

				_, err := mock.userUC.GetRoleByID(1)
				assert.Equal(t, nil, err)
			},
		},
	})
}
