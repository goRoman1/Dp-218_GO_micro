package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
)

// UserService - structure for implementing user service
type UserService struct {
	repoUser repositories.UserRepo
	repoRole repositories.RoleRepo
}

// NewUserService - initialization of UserService
func NewUserService(repoUser repositories.UserRepo, repoRole repositories.RoleRepo) *UserService {
	return &UserService{repoUser: repoUser, repoRole: repoRole}
}

// GetAllUsers - get all system users
func (ser *UserService) GetAllUsers() (*models.UserList, error) {
	return ser.repoUser.GetAllUsers()
}

// AddUser - create new system user
func (ser *UserService) AddUser(user *models.User) error {
	return ser.repoUser.AddUser(user)
}

// GetUserByID - get user information by its ID
func (ser *UserService) GetUserByID(userID int) (models.User, error) {
	return ser.repoUser.GetUserByID(userID)
}

// DeleteUser - delete given user by its ID
func (ser *UserService) DeleteUser(userID int) error {
	return ser.repoUser.DeleteUser(userID)
}

// UpdateUser - update user information by its ID
func (ser *UserService) UpdateUser(userID int, userData models.User) (models.User, error) {
	return ser.repoUser.UpdateUser(userID, userData)
}

// FindUsersByLoginNameSurname - find system user by given login(email), name or surname
func (ser *UserService) FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error) {
	return ser.repoUser.FindUsersByLoginNameSurname(whatToFind)
}

// GetAllRoles - get all roles information
func (ser *UserService) GetAllRoles() (*models.RoleList, error) {
	return ser.repoRole.GetAllRoles()
}

// GetRoleByID - get role information by role ID
func (ser *UserService) GetRoleByID(roleID int) (models.Role, error) {
	return ser.repoRole.GetRoleByID(roleID)
}

// ChangeUsersBlockStatus - change user blocked status to the opposite (true->false, false->true)
func (ser *UserService) ChangeUsersBlockStatus(userID int) error {
	user, err := ser.repoUser.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.IsBlocked = !user.IsBlocked
	_, err = ser.repoUser.UpdateUser(userID, user)
	return err
}
