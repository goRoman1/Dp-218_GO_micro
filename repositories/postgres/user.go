package postgres

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"context"
	"fmt"
	"time"
)

// UserRepoDB - struct representing user repository
type UserRepoDB struct {
	db repositories.AnyDatabase
}

// NewUserRepoDB - user repo initialization
func NewUserRepoDB(db repositories.AnyDatabase) *UserRepoDB {
	return &UserRepoDB{db}
}

// GetAllUsers - get list of all system users from the DB
func (urdb *UserRepoDB) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}

	roles, err := urdb.GetAllRoles()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT 
		id, login_email, is_blocked, user_name, user_surname, created_at, role_id 
		FROM users 
		ORDER BY id DESC;`
	rows, err := urdb.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var user models.User
		var roleID int
		err := rows.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
			&user.UserName, &user.UserSurname, &user.CreatedAt, &roleID)
		if err != nil {
			return list, err
		}

		user.Role, err = FindRoleInTheList(roles, roleID)
		if err != nil {
			return list, err
		}

		list.Users = append(list.Users, user)
	}
	return list, nil
}

// AddUser - create user record in the DB based on given entity
func (urdb *UserRepoDB) AddUser(user *models.User) error {
	var id int
	var createdAt time.Time
	querySQL := `INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id, password_hash) 
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at;`
	err := urdb.db.QueryResultRow(context.Background(), querySQL,
		user.LoginEmail, user.IsBlocked, user.UserName, user.UserSurname, user.Role.ID, user.Password).
		Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = createdAt
	return nil
}

// GetUserByID - get user entity from the DB by given user ID
func (urdb *UserRepoDB) GetUserByID(userID int) (models.User, error) {
	user := models.User{}

	querySQL := `SELECT 
		id, login_email, is_blocked, user_name, user_surname, created_at, role_id
		FROM users 
		WHERE id = $1;`
	row := urdb.db.QueryResultRow(context.Background(), querySQL, userID)

	var roleID int
	err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &roleID)
	if err != nil {
		return models.User{}, err
	}
	user.Role, err = urdb.GetRoleByID(roleID)

	return user, err
}

// GetUserByEmail - get user entity from the DB by given user email
func (urdb *UserRepoDB) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	querySQL := `SELECT 
		id, login_email, is_blocked, user_name, user_surname, created_at, role_id, password_hash 
		FROM users 
		WHERE login_email = $1;`
	row := urdb.db.QueryResultRow(context.Background(), querySQL, email)

	var roleID int
	err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &roleID, &user.Password)

	if err != nil {
		return models.User{}, err
	}
	user.Role, err = urdb.GetRoleByID(roleID)

	return user, err
}

// DeleteUser - delete user with given ID from the DB
func (urdb *UserRepoDB) DeleteUser(userID int) error {
	querySQL := `DELETE FROM users WHERE id = $1;`
	_, err := urdb.db.QueryExec(context.Background(), querySQL, userID)
	return err
}

// UpdateUser - update user with given ID in the DB based on given user entity
func (urdb *UserRepoDB) UpdateUser(userID int, userData models.User) (models.User, error) {
	user := models.User{}
	querySQL := `UPDATE users 
		SET login_email=$1, is_blocked=$2, user_name=$3, user_surname=$4, role_id=$5 
		WHERE id=$6 
		RETURNING id, created_at, login_email, is_blocked, user_name, user_surname, role_id;`
	var roleID int
	err := urdb.db.QueryResultRow(context.Background(), querySQL,
		userData.LoginEmail, userData.IsBlocked, userData.UserName,
		userData.UserSurname, userData.Role.ID, userID).
		Scan(&user.ID, &user.CreatedAt, &user.LoginEmail, &user.IsBlocked, &user.UserName, &user.UserSurname, &roleID)
	if err != nil {
		return user, err
	}
	user.Role, err = urdb.GetRoleByID(roleID)
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindUsersByLoginNameSurname - find list of users having whatToFind string in login, name or surname in the DB
func (urdb *UserRepoDB) FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error) {
	list := &models.UserList{}

	roles, err := urdb.GetAllRoles()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT id, login_email, is_blocked, user_name, user_surname, created_at, role_id FROM users 
		WHERE LOWER(login_email) LIKE LOWER($1) 
			OR LOWER(user_name) LIKE LOWER($1) 
			OR LOWER(user_surname) LIKE LOWER($1) 
		ORDER BY id DESC;`
	rows, err := urdb.db.QueryResult(context.Background(), querySQL, whatToFind+"%")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var user models.User
		var roleID int
		err := rows.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
			&user.UserName, &user.UserSurname, &user.CreatedAt, &roleID)
		if err != nil {
			return list, err
		}

		user.Role, err = FindRoleInTheList(roles, roleID)
		if err != nil {
			return list, err
		}

		list.Users = append(list.Users, user)
	}
	return list, nil
}

// GetAllRoles - get list of all roles from the DB
func (urdb *UserRepoDB) GetAllRoles() (*models.RoleList, error) {
	list := &models.RoleList{}
	querySQL := `SELECT * FROM roles ORDER BY id DESC;`
	rows, err := urdb.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.IsAdmin, &role.IsUser, &role.IsSupplier)
		if err != nil {
			return list, err
		}
		list.Roles = append(list.Roles, role)
	}
	return list, nil
}

// GetRoleByID - get role from the DB by given role ID
func (urdb *UserRepoDB) GetRoleByID(roleId int) (models.Role, error) {
	role := models.Role{}
	querySQL := `SELECT * FROM roles WHERE id = $1;`
	row := urdb.db.QueryResultRow(context.Background(), querySQL, roleId)
	err := row.Scan(&role.ID, &role.Name, &role.IsAdmin, &role.IsUser, &role.IsSupplier)
	return role, err
}

// FindRoleInTheList - find role in given role list by role ID
func FindRoleInTheList(roles *models.RoleList, roleID int) (models.Role, error) {
	for _, v := range roles.Roles {
		if v.ID == roleID {
			return v, nil
		}
	}
	return models.Role{}, fmt.Errorf("not found role id=%d", roleID)
}
