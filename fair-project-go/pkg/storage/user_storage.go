package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sadeq/fair-project-go/pkg/models"
)

const (
	usersFilename = "users.json"
)

var (
	usersMutex sync.RWMutex
)

// SaveUsers saves the users to data/users.json
func SaveUsers(users map[string]models.User) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	if err := EnsureBaseDir(); err != nil {
		return fmt.Errorf("failed to ensure base data directory: %w", err)
	}

	filePath := filepath.Join(baseDataDir, usersFilename)
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write users to file '%s': %w", filePath, err)
	}
	return nil
}

// LoadUsers loads users from data/users.json
func LoadUsers() (map[string]models.User, error) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	if err := EnsureBaseDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure base data directory: %w", err)
	}

	filePath := filepath.Join(baseDataDir, usersFilename)

	// If the file doesn't exist, return an empty map
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return make(map[string]models.User), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file '%s': %w", filePath, err)
	}

	var users map[string]models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users from JSON file '%s': %w", filePath, err)
	}
	return users, nil
}

// GetUser gets a user by username
func GetUser(username string) (models.User, bool, error) {
	users, err := LoadUsers()
	if err != nil {
		return models.User{}, false, err
	}

	user, exists := users[username]
	return user, exists, nil
}

// SaveUser saves a user
func SaveUser(user models.User) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	// Update the timestamps
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.UpdatedAt
	}

	users[user.Username] = user
	return SaveUsers(users)
}

// DeleteUser deletes a user
func DeleteUser(username string) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	delete(users, username)
	return SaveUsers(users)
}

// InitializeDefaultUsers creates default users if no users exist
func InitializeDefaultUsers() error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	// If users already exist, do nothing
	if len(users) > 0 {
		return nil
	}

	// Create a default superadmin user
	// In a real application, you would use a secure password hashing algorithm
	// For this example, we'll use a simple hash (not secure for production)
	superadmin := models.User{
		Username:     "superadmin",
		PasswordHash: "$2a$10$rNQZQSJQhECnhMOuFITnA.FZbCVHvdAZQTcAOVYxuRjEaY4ykFpS2", // "password"
		Roles:        []models.Role{models.RoleSuperAdmin},
		Enabled:      true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	users[superadmin.Username] = superadmin
	return SaveUsers(users)
}
