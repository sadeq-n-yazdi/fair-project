package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// AssignmentRepository implements the repository.AssignmentRepository interface for PostgreSQL storage
type AssignmentRepository struct {
	*BaseRepository
}

// NewAssignmentRepository creates a new PostgreSQL assignment repository
func NewAssignmentRepository(db *sql.DB) repository.AssignmentRepository {
	return &AssignmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Initialize initializes the repository by creating the necessary tables
func (r *AssignmentRepository) Initialize(ctx context.Context) error {
	// Call the base Initialize method
	if err := r.BaseRepository.Initialize(ctx); err != nil {
		return err
	}

	// Create the assignments table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS assignments (
		class_name TEXT NOT NULL,
		assignment_id TEXT NOT NULL,
		results TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (class_name, assignment_id)
	)
	`
	_, err := r.ExecContext(ctx, query)
	return err
}

// SaveAssignmentResults saves assignment results for a class/term
func (r *AssignmentRepository) SaveAssignmentResults(ctx context.Context, className string, assignmentID string, results models.AssignmentOutput) error {
	// Convert results to JSON
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal assignment results: %w", err)
	}

	// Check if the assignment already exists
	var exists bool
	query := `SELECT 1 FROM assignments WHERE class_name = $1 AND assignment_id = $2`
	row := r.QueryRowContext(ctx, query, className, assignmentID)
	err = row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if assignment exists: %w", err)
	}

	// Insert or update the assignment
	if err == sql.ErrNoRows {
		// Insert
		query = `INSERT INTO assignments (class_name, assignment_id, results, created_at) VALUES ($1, $2, $3, NOW())`
		_, err = r.ExecContext(ctx, query, className, assignmentID, string(resultsJSON))
	} else {
		// Update
		query = `UPDATE assignments SET results = $1, created_at = NOW() WHERE class_name = $2 AND assignment_id = $3`
		_, err = r.ExecContext(ctx, query, string(resultsJSON), className, assignmentID)
	}

	if err != nil {
		return fmt.Errorf("failed to save assignment results: %w", err)
	}

	return nil
}

// ListAssignmentResults lists all assignment results for a class/term
func (r *AssignmentRepository) ListAssignmentResults(ctx context.Context, className string) ([]string, error) {
	query := `SELECT assignment_id FROM assignments WHERE class_name = $1 ORDER BY created_at DESC`
	rows, err := r.QueryContext(ctx, query, className)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignment results: %w", err)
	}
	defer rows.Close()

	var assignmentIDs []string
	for rows.Next() {
		var assignmentID string
		if err := rows.Scan(&assignmentID); err != nil {
			return nil, fmt.Errorf("failed to scan assignment ID: %w", err)
		}
		assignmentIDs = append(assignmentIDs, assignmentID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assignment IDs: %w", err)
	}

	return assignmentIDs, nil
}

// LoadAssignmentResult loads an assignment result for a class/term
func (r *AssignmentRepository) LoadAssignmentResult(ctx context.Context, className string, assignmentID string) (models.AssignmentOutput, error) {
	query := `SELECT results FROM assignments WHERE class_name = $1 AND assignment_id = $2`
	row := r.QueryRowContext(ctx, query, className, assignmentID)

	var resultsJSON string
	err := row.Scan(&resultsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.AssignmentOutput{}, fmt.Errorf("assignment not found")
		}
		return models.AssignmentOutput{}, fmt.Errorf("failed to load assignment result: %w", err)
	}

	var results models.AssignmentOutput
	if err := json.Unmarshal([]byte(resultsJSON), &results); err != nil {
		return models.AssignmentOutput{}, fmt.Errorf("failed to unmarshal assignment results: %w", err)
	}

	return results, nil
}
