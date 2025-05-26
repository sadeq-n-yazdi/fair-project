package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// ClassRepository implements the repository.ClassRepository interface for PostgreSQL storage
type ClassRepository struct {
	*BaseRepository
}

// NewClassRepository creates a new PostgreSQL class repository
func NewClassRepository(db *sql.DB) repository.ClassRepository {
	return &ClassRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Initialize initializes the repository by creating the necessary tables
func (r *ClassRepository) Initialize(ctx context.Context) error {
	// Call the base Initialize method
	if err := r.BaseRepository.Initialize(ctx); err != nil {
		return err
	}

	// Create the classes table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS classes (
		name TEXT PRIMARY KEY,
		created_at TIMESTAMP NOT NULL
	)
	`
	_, err := r.ExecContext(ctx, query)
	return err
}

// CreateClassTerm creates a new class/term
func (r *ClassRepository) CreateClassTerm(ctx context.Context, className string) error {
	// Check if the class already exists
	var exists bool
	query := `SELECT 1 FROM classes WHERE name = $1`
	row := r.QueryRowContext(ctx, query, className)
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if class exists: %w", err)
	}

	// If the class doesn't exist, create it
	if err == sql.ErrNoRows {
		query = `INSERT INTO classes (name, created_at) VALUES ($1, NOW())`
		_, err = r.ExecContext(ctx, query, className)
		if err != nil {
			return fmt.Errorf("failed to create class: %w", err)
		}
	}

	return nil
}

// ListClassTerms lists all class/terms
func (r *ClassRepository) ListClassTerms(ctx context.Context) ([]string, error) {
	query := `SELECT name FROM classes ORDER BY name`
	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list classes: %w", err)
	}
	defer rows.Close()

	var classes []string
	for rows.Next() {
		var className string
		if err := rows.Scan(&className); err != nil {
			return nil, fmt.Errorf("failed to scan class name: %w", err)
		}
		classes = append(classes, className)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating classes: %w", err)
	}

	return classes, nil
}
