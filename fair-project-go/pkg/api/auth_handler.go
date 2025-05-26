package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sadeq/fair-project-go/pkg/auth"
	"github.com/sadeq/fair-project-go/pkg/errors"
	authmiddleware "github.com/sadeq/fair-project-go/pkg/middleware/auth"
	"github.com/sadeq/fair-project-go/pkg/models"
)

// LoginHandler handles POST requests to /auth/login
// It authenticates a user and returns a JWT token
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only POST method is allowed"))
		return
	}

	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body"))
		return
	}
	defer r.Body.Close()

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Get the user from the repository
	ctx := r.Context()
	user, exists, err := userRepo.GetUser(ctx, credentials.Username)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to get user"))
		return
	}

	if !exists {
		HandleError(w, errors.New(ErrInvalidCredentials, "Invalid username or password"))
		return
	}

	// Check if the user is enabled
	if !user.Enabled {
		HandleError(w, errors.New(ErrInvalidCredentials, "User is disabled"))
		return
	}

	// Check if the password is correct
	if !auth.CheckPassword(credentials.Password, user.PasswordHash) {
		HandleError(w, errors.New(ErrInvalidCredentials, "Invalid username or password"))
		return
	}

	// Generate a JWT token
	token, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to generate token"))
		return
	}

	// Return the token
	response := models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}

	RespondJSON(w, http.StatusOK, response)
}

// CreateUserHandler handles POST requests to /auth/users
// It creates a new user (admin and superadmin only)
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only POST method is allowed"))
		return
	}

	// Check if the user has permission to create users
	currentUser, ok := authmiddleware.GetUserFromContext(r)
	if !ok || !currentUser.CanManageUsers() {
		HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions"))
		return
	}

	var newUser models.NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body"))
		return
	}
	defer r.Body.Close()

	// Validate the request
	if newUser.Username == "" {
		HandleError(w, errors.New(errors.ErrInvalidRequestBody, "Username is required"))
		return
	}

	if newUser.Password == "" {
		HandleError(w, errors.New(errors.ErrInvalidRequestBody, "Password is required"))
		return
	}

	if len(newUser.Roles) == 0 {
		HandleError(w, errors.New(ErrInvalidRoles, "At least one role is required"))
		return
	}

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Check if the user already exists
	ctx := r.Context()
	_, exists, err := userRepo.GetUser(ctx, newUser.Username)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to check if user exists"))
		return
	}

	if exists {
		HandleError(w, errors.New(ErrUserAlreadyExists, "User already exists"))
		return
	}

	// Check if the current user has permission to assign the requested roles
	// Only superadmins can create other superadmins
	if !currentUser.CanManageRoles() {
		for _, role := range newUser.Roles {
			if role == models.RoleSuperAdmin {
				HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions to assign superadmin role"))
				return
			}
		}
	}

	// Hash the password
	passwordHash, err := auth.HashPassword(newUser.Password)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to hash password"))
		return
	}

	// Create the user
	user := models.User{
		Username:     newUser.Username,
		PasswordHash: passwordHash,
		Roles:        newUser.Roles,
		Enabled:      true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the user
	if err := userRepo.SaveUser(ctx, user); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to save user"))
		return
	}

	// Return the user (without password hash)
	RespondJSON(w, http.StatusCreated, user)
}

// UpdateUserHandler handles PUT requests to /auth/users/{username}
// It updates a user's roles and enabled status (admin and superadmin only)
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request, username string) {
	if r.Method != http.MethodPut {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only PUT method is allowed"))
		return
	}

	// Check if the user has permission to update users
	currentUser, ok := authmiddleware.GetUserFromContext(r)
	if !ok || !currentUser.CanManageUsers() {
		HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions"))
		return
	}

	var updateRequest models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body"))
		return
	}
	defer r.Body.Close()

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Get the user to update
	ctx := r.Context()
	user, exists, err := userRepo.GetUser(ctx, username)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to get user"))
		return
	}

	if !exists {
		HandleError(w, errors.New(ErrUserNotFound, "User not found"))
		return
	}

	// Check if the current user has permission to update the user's roles
	// Only superadmins can update roles to include superadmin
	if updateRequest.Roles != nil {
		if !currentUser.CanManageRoles() {
			for _, role := range updateRequest.Roles {
				if role == models.RoleSuperAdmin {
					HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions to assign superadmin role"))
					return
				}
			}
		}

		// Ensure at least one role
		if len(updateRequest.Roles) == 0 {
			HandleError(w, errors.New(ErrInvalidRoles, "At least one role is required"))
			return
		}

		user.Roles = updateRequest.Roles
	}

	// Update enabled status if provided
	if updateRequest.Enabled != nil {
		user.Enabled = *updateRequest.Enabled
	}

	// Update the timestamp
	user.UpdatedAt = time.Now()

	// Save the user
	if err := userRepo.SaveUser(ctx, user); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to save user"))
		return
	}

	// Return the updated user
	RespondJSON(w, http.StatusOK, user)
}

// GetUsersHandler handles GET requests to /auth/users
// It returns a list of all users (admin and superadmin only)
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed"))
		return
	}

	// Check if the user has permission to list users
	currentUser, ok := authmiddleware.GetUserFromContext(r)
	if !ok || !currentUser.CanManageUsers() {
		HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions"))
		return
	}

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Get all users
	ctx := r.Context()
	users, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to load users"))
		return
	}

	// Convert map to slice for easier consumption by clients
	usersList := make([]models.User, 0, len(users))
	for _, user := range users {
		usersList = append(usersList, user)
	}

	RespondJSON(w, http.StatusOK, usersList)
}

// GetUserHandler handles GET requests to /auth/users/{username}
// It returns a specific user (admin and superadmin only, or the user themselves)
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request, username string) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed"))
		return
	}

	// Check if the user has permission to view this user
	currentUser, ok := authmiddleware.GetUserFromContext(r)
	if !ok || (!currentUser.CanManageUsers() && currentUser.Username != username) {
		HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions"))
		return
	}

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Get the user
	ctx := r.Context()
	user, exists, err := userRepo.GetUser(ctx, username)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to get user"))
		return
	}

	if !exists {
		HandleError(w, errors.New(ErrUserNotFound, "User not found"))
		return
	}

	RespondJSON(w, http.StatusOK, user)
}

// DeleteUserHandler handles DELETE requests to /auth/users/{username}
// It deletes a user (superadmin only)
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request, username string) {
	if r.Method != http.MethodDelete {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only DELETE method is allowed"))
		return
	}

	// Check if the user has permission to delete users
	currentUser, ok := authmiddleware.GetUserFromContext(r)
	if !ok || !currentUser.CanManageRoles() {
		HandleError(w, errors.New(authmiddleware.ErrForbidden, "Insufficient permissions"))
		return
	}

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Check if the user exists
	ctx := r.Context()
	_, exists, err := userRepo.GetUser(ctx, username)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to check if user exists"))
		return
	}

	if !exists {
		HandleError(w, errors.New(ErrUserNotFound, "User not found"))
		return
	}

	// Don't allow deleting yourself
	if currentUser.Username == username {
		HandleError(w, errors.New(errors.ErrInvalidRequestBody, "Cannot delete yourself"))
		return
	}

	// Delete the user
	if err := userRepo.DeleteUser(ctx, username); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to delete user"))
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// WhoAmIHandler handles GET requests to /user/whoami
// It returns the current user's information if authenticated
// If not authenticated, it returns a JSON response indicating that
func (h *Handler) WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed"))
		return
	}

	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")

	// If no Authorization header, return a response indicating not authenticated
	if authHeader == "" {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "No JWT token provided",
		})
		return
	}

	// Check if the Authorization header has the correct format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid authorization header format",
		})
		return
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate the token
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid or expired token",
		})
		return
	}

	// Get the user repository from the storage manager
	userRepo := h.storageManager.UserRepository()

	// Get the user from the repository
	ctx := r.Context()
	username := claims.Username
	user, exists, err := userRepo.GetUser(ctx, username)
	if err != nil || !exists {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "User not found",
		})
		return
	}

	if !user.Enabled {
		RespondJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "User is disabled",
		})
		return
	}

	// Return the user information
	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"username":      user.Username,
		"status":        user.Enabled,
		"roles":         user.Roles,
	})
}

// AuthResourceHandler handles requests to /auth/{resource}
// It routes to the appropriate handler based on the resource
func (h *Handler) AuthResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/auth/login" {
		h.LoginHandler(w, r)
		return
	}

	if path == "/user/whoami" {
		h.WhoAmIHandler(w, r)
		return
	}

	if path == "/auth/users" {
		if r.Method == http.MethodGet {
			h.GetUsersHandler(w, r)
		} else if r.Method == http.MethodPost {
			h.CreateUserHandler(w, r)
		} else {
			HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Method not allowed"))
		}
		return
	}

	// Handle /auth/users/{username}
	if len(path) > 12 && path[:12] == "/auth/users/" {
		username := path[12:]
		if r.Method == http.MethodGet {
			h.GetUserHandler(w, r, username)
		} else if r.Method == http.MethodPut {
			h.UpdateUserHandler(w, r, username)
		} else if r.Method == http.MethodDelete {
			h.DeleteUserHandler(w, r, username)
		} else {
			HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Method not allowed"))
		}
		return
	}

	// If we get here, the resource was not found
	HandleError(w, errors.New(errors.ErrResourceNotFound, "Resource not found"))
}
