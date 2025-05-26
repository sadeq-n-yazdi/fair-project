package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

const (
	usersFilename = "users.json"
)

// UserRepository implements the repository.UserRepository interface for file-based storage
type UserRepository struct {
	*BaseRepository
	mu sync.RWMutex
}

// NewUserRepository creates a new file-based user repository
func NewUserRepository(baseDir string) repository.UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(baseDir),
	}
}

// GetUser gets a user by username
func (r *UserRepository) GetUser(ctx context.Context, username string) (models.User, bool, error) {
	users, err := r.GetAllUsers(ctx)
	if err != nil {
		return models.User{}, false, err
	}

	user, exists := users[username]
	return user, exists, nil
}

// SaveUser saves a user
func (r *UserRepository) SaveUser(ctx context.Context, user models.User) error {
	users, err := r.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	// Update the timestamps
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.UpdatedAt
	}

	users[user.Username] = user
	return r.saveUsers(ctx, users)
}

// DeleteUser deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	users, err := r.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	delete(users, username)
	return r.saveUsers(ctx, users)
}

// GetAllUsers gets all users
func (r *UserRepository) GetAllUsers(ctx context.Context) (map[string]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := r.EnsureBaseDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure base data directory: %w", err)
	}

	filePath := filepath.Join(r.baseDir, usersFilename)

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

// InitializeDefaultUsers creates default users if no users exist
func (r *UserRepository) InitializeDefaultUsers(ctx context.Context) error {
	users, err := r.GetAllUsers(ctx)
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
	return r.saveUsers(ctx, users)
}

// saveUsers saves the users to the users file
func (r *UserRepository) saveUsers(ctx context.Context, users map[string]models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.EnsureBaseDir(); err != nil {
		return fmt.Errorf("failed to ensure base data directory: %w", err)
	}

	filePath := filepath.Join(r.baseDir, usersFilename)
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write users to file '%s': %w", filePath, err)
	}
	return nil
}
