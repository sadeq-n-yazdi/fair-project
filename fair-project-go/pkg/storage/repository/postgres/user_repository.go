package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// UserRepository implements the repository.UserRepository interface for PostgreSQL storage
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Initialize initializes the repository by creating the necessary tables
func (r *UserRepository) Initialize(ctx context.Context) error {
	// Call the base Initialize method
	if err := r.BaseRepository.Initialize(ctx); err != nil {
		return err
	}

	// Create the users table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		roles TEXT NOT NULL,
		enabled BOOLEAN NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)
	`
	_, err := r.ExecContext(ctx, query)
	return err
}

// GetUser gets a user by username
func (r *UserRepository) GetUser(ctx context.Context, username string) (models.User, bool, error) {
	query := `
	SELECT username, password_hash, roles, enabled, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	row := r.QueryRowContext(ctx, query, username)

	var user models.User
	var rolesJSON string
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Username, &user.PasswordHash, &rolesJSON, &user.Enabled, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, false, nil
		}
		return models.User{}, false, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse the roles JSON
	if err := json.Unmarshal([]byte(rolesJSON), &user.Roles); err != nil {
		return models.User{}, false, fmt.Errorf("failed to unmarshal roles: %w", err)
	}

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, true, nil
}

// SaveUser saves a user
func (r *UserRepository) SaveUser(ctx context.Context, user models.User) error {
	// Update the timestamps
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.UpdatedAt
	}

	// Convert the roles to JSON
	rolesJSON, err := json.Marshal(user.Roles)
	if err != nil {
		return fmt.Errorf("failed to marshal roles: %w", err)
	}

	// Check if the user exists
	var exists bool
	query := `SELECT 1 FROM users WHERE username = $1`
	row := r.QueryRowContext(ctx, query, user.Username)
	err = row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	// Insert or update the user
	if err == sql.ErrNoRows {
		// Insert
		query = `
		INSERT INTO users (username, password_hash, roles, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = r.ExecContext(ctx, query, user.Username, user.PasswordHash, string(rolesJSON),
			user.Enabled, user.CreatedAt, user.UpdatedAt)
	} else {
		// Update
		query = `
		UPDATE users
		SET password_hash = $1, roles = $2, enabled = $3, updated_at = $4
		WHERE username = $5
		`
		_, err = r.ExecContext(ctx, query, user.PasswordHash, string(rolesJSON),
			user.Enabled, user.UpdatedAt, user.Username)
	}

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = $1`
	_, err := r.ExecContext(ctx, query, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetAllUsers gets all users
func (r *UserRepository) GetAllUsers(ctx context.Context) (map[string]models.User, error) {
	query := `
	SELECT username, password_hash, roles, enabled, created_at, updated_at
	FROM users
	`
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := make(map[string]models.User)
	for rows.Next() {
		var user models.User
		var rolesJSON string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&user.Username, &user.PasswordHash, &rolesJSON, &user.Enabled, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		// Parse the roles JSON
		if err := json.Unmarshal([]byte(rolesJSON), &user.Roles); err != nil {
			return nil, fmt.Errorf("failed to unmarshal roles: %w", err)
		}

		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt

		users[user.Username] = user
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
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
	superadmin := models.User{
		Username:     "superadmin",
		PasswordHash: "$2a$10$rNQZQSJQhECnhMOuFITnA.FZbCVHvdAZQTcAOVYxuRjEaY4ykFpS2", // "password"
		Roles:        []models.Role{models.RoleSuperAdmin},
		Enabled:      true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return r.SaveUser(ctx, superadmin)
}
