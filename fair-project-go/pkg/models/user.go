package models

import (
	"time"
)

// Role represents a user role in the system
type Role string

const (
	// RoleSuperAdmin is the highest level role with full access
	RoleSuperAdmin Role = "superadmin"
	// RoleAdmin can manage classes and users
	RoleAdmin Role = "admin"
	// RoleTeacher can view and update existing projects
	RoleTeacher Role = "teacher"
	// RoleStudent can only view data
	RoleStudent Role = "student"
)

// User represents a user in the system
type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash,omitempty"` // Include in storage but can be omitted in responses
	Roles        []Role    `json:"roles"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (u *User) HasAnyRole(roles ...Role) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

// CanManageUsers returns true if the user can manage users (admin or superadmin)
func (u *User) CanManageUsers() bool {
	return u.HasAnyRole(RoleAdmin, RoleSuperAdmin)
}

// CanManageRoles returns true if the user can manage roles (superadmin only)
func (u *User) CanManageRoles() bool {
	return u.HasRole(RoleSuperAdmin)
}

// CanDefineClasses returns true if the user can define classes (admin or superadmin)
func (u *User) CanDefineClasses() bool {
	return u.HasAnyRole(RoleAdmin, RoleSuperAdmin)
}

// CanUpdateProjects returns true if the user can update projects (teacher, admin, or superadmin)
func (u *User) CanUpdateProjects() bool {
	return u.HasAnyRole(RoleTeacher, RoleAdmin, RoleSuperAdmin)
}

// CanViewData returns true if the user can view data (any role)
func (u *User) CanViewData() bool {
	return len(u.Roles) > 0
}

// UserCredentials represents the credentials sent during login
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUserRequest represents the request to create a new user
type NewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    []Role `json:"roles"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Roles   []Role `json:"roles,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

// TokenResponse represents the response after successful authentication
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}
