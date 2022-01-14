package models

import (
	"time"
)

// Role - entity for user roles
type Role struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsAdmin    bool   `json:"is_admin"`
	IsUser     bool   `json:"is_user"`
	IsSupplier bool   `json:"is_supplier"`
}

// RoleList - struct with list of roles
type RoleList struct {
	Roles []Role `json:"roles"`
}

// User - entity representing user in the system
type User struct {
	ID          int       `json:"id"`
	LoginEmail  string    `json:"login_email"`
	IsBlocked   bool      `json:"is_blocked"`
	UserName    string    `json:"user_name"`
	UserSurname string    `json:"user_surname"`
	CreatedAt   time.Time `json:"created_at"`
	Role        Role      `json:"role"`
	Password    string    `json:"password"`
}

// UserList - struct for list of users
type UserList struct {
	Users []User `json:"users"`
}
