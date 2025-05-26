package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// ProjectRepository implements the repository.ProjectRepository interface for PostgreSQL storage
type ProjectRepository struct {
	*BaseRepository
}

// NewProjectRepository creates a new PostgreSQL project repository
func NewProjectRepository(db *sql.DB) repository.ProjectRepository {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Initialize initializes the repository by creating the necessary tables
func (r *ProjectRepository) Initialize(ctx context.Context) error {
	// Call the base Initialize method
	if err := r.BaseRepository.Initialize(ctx); err != nil {
		return err
	}

	// Create the projects table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS projects (
		class_name TEXT NOT NULL,
		projects TEXT NOT NULL,
		PRIMARY KEY (class_name)
	)
	`
	_, err := r.ExecContext(ctx, query)
	return err
}

// SaveProjects saves projects for a class/term
func (r *ProjectRepository) SaveProjects(ctx context.Context, className string, projects []string) error {
	// Convert projects to JSON
	projectsJSON, err := json.Marshal(projects)
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	// Check if the class already has projects
	var exists bool
	query := `SELECT 1 FROM projects WHERE class_name = $1`
	row := r.QueryRowContext(ctx, query, className)
	err = row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if projects exist: %w", err)
	}

	// Insert or update the projects
	if err == sql.ErrNoRows {
		// Insert
		query = `INSERT INTO projects (class_name, projects) VALUES ($1, $2)`
		_, err = r.ExecContext(ctx, query, className, string(projectsJSON))
	} else {
		// Update
		query = `UPDATE projects SET projects = $1 WHERE class_name = $2`
		_, err = r.ExecContext(ctx, query, string(projectsJSON), className)
	}

	if err != nil {
		return fmt.Errorf("failed to save projects: %w", err)
	}

	return nil
}

// LoadProjects loads projects for a class/term
func (r *ProjectRepository) LoadProjects(ctx context.Context, className string) ([]string, error) {
	query := `SELECT projects FROM projects WHERE class_name = $1`
	row := r.QueryRowContext(ctx, query, className)

	var projectsJSON string
	err := row.Scan(&projectsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load projects: %w", err)
	}

	var projects []string
	if err := json.Unmarshal([]byte(projectsJSON), &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return projects, nil
}

// ReadProjectsFromFile reads projects from a file
func (r *ProjectRepository) ReadProjectsFromFile(ctx context.Context, fileName string) ([]string, error) {
	// This method doesn't interact with the database, so it's the same for all storage types
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects file: %w", err)
	}

	var projects []string
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return projects, nil
}
