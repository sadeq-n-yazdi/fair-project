package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// StudentRepository implements the repository.StudentRepository interface for PostgreSQL storage
type StudentRepository struct {
	*BaseRepository
}

// NewStudentRepository creates a new PostgreSQL student repository
func NewStudentRepository(db *sql.DB) repository.StudentRepository {
	return &StudentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Initialize initializes the repository by creating the necessary tables
func (r *StudentRepository) Initialize(ctx context.Context) error {
	// Call the base Initialize method
	if err := r.BaseRepository.Initialize(ctx); err != nil {
		return err
	}

	// Create the students table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS students (
		class_name TEXT NOT NULL,
		students TEXT NOT NULL,
		PRIMARY KEY (class_name)
	)
	`
	_, err := r.ExecContext(ctx, query)
	return err
}

// SaveStudents saves students for a class/term
func (r *StudentRepository) SaveStudents(ctx context.Context, className string, students []models.StudentInput) error {
	// Convert students to JSON
	studentsJSON, err := json.Marshal(students)
	if err != nil {
		return fmt.Errorf("failed to marshal students: %w", err)
	}

	// Check if the class already has students
	var exists bool
	query := `SELECT 1 FROM students WHERE class_name = $1`
	row := r.QueryRowContext(ctx, query, className)
	err = row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if students exist: %w", err)
	}

	// Insert or update the students
	if err == sql.ErrNoRows {
		// Insert
		query = `INSERT INTO students (class_name, students) VALUES ($1, $2)`
		_, err = r.ExecContext(ctx, query, className, string(studentsJSON))
	} else {
		// Update
		query = `UPDATE students SET students = $1 WHERE class_name = $2`
		_, err = r.ExecContext(ctx, query, string(studentsJSON), className)
	}

	if err != nil {
		return fmt.Errorf("failed to save students: %w", err)
	}

	return nil
}

// LoadStudents loads students for a class/term
func (r *StudentRepository) LoadStudents(ctx context.Context, className string) ([]models.StudentInput, error) {
	query := `SELECT students FROM students WHERE class_name = $1`
	row := r.QueryRowContext(ctx, query, className)

	var studentsJSON string
	err := row.Scan(&studentsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load students: %w", err)
	}

	var students []models.StudentInput
	if err := json.Unmarshal([]byte(studentsJSON), &students); err != nil {
		return nil, fmt.Errorf("failed to unmarshal students: %w", err)
	}

	return students, nil
}

// ReadStudentsFromFile reads students from a file
func (r *StudentRepository) ReadStudentsFromFile(ctx context.Context, fileName string) ([]models.Student, error) {
	// This method doesn't interact with the database, so it's the same for all storage types
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read students file: %w", err)
	}

	var students []models.Student
	if err := json.Unmarshal(data, &students); err != nil {
		return nil, fmt.Errorf("failed to unmarshal students: %w", err)
	}

	return students, nil
}
