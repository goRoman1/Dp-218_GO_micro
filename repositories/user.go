//go:generate mockgen -source=user.go -destination=../repositories/mock/mock_user.go -package=mock
package repositories

import (
	"Dp218GO/models"
	"context"
)

// UserRepo - interface for user repository
type UserRepo interface {
	GetAllUsers() (*models.UserList, error)
	GetUserByID(userID int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userID int, userData models.User) (models.User, error)
	DeleteUser(userID int) error
	FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error)
}

// RoleRepo - interface for role repository
type RoleRepo interface {
	GetAllRoles() (*models.RoleList, error)
	GetRoleByID(roleID int) (models.Role, error)
}

// AuthRepo - interface for authorization repository
type AuthRepo interface {
	GetUserByEmail(context.Context, string) (models.User, error)
}
