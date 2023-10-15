package internal

import "github.com/dark-enstein/port/auth"

type Internal interface{}

// Repository defines a data structure that maps to a table in the db
type Repository interface {
	Create() ([]byte, error)
}

type User interface {
	GetPermissions() auth.PermissionSet
	GetRoles() auth.RoleSet
	IsValid()
}

type Agent struct {
	Name  Name         `json:"name"`
	ID    []byte       `json:"id"`
	Roles auth.RoleSet `json:"roles,omitempty"`
}

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type DateOfBirth struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

func (a *Agent) GetPermissions() auth.PermissionSet {
	// db calls
	return auth.PermissionSet{}
}

func (a *Agent) GetRoles() auth.RoleSet {
	// db calls
	return auth.RoleSet{}
}
