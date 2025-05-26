package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

const (
	studentsFilename = "students.json"
)

// StudentRepository implements the repository.StudentRepository interface for file-based storage
type StudentRepository struct {
	*BaseRepository
}

// NewStudentRepository creates a new file-based student repository
func NewStudentRepository(baseDir string) repository.StudentRepository {
	return &StudentRepository{
		BaseRepository: NewBaseRepository(baseDir),
	}
}

// SaveStudents saves students for a class/term
func (r *StudentRepository) SaveStudents(ctx context.Context, className string, students []models.StudentInput) error {
	if !r.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}

	if err := r.EnsureClassTermDir(className); err != nil {
		return err
	}

	filePath := filepath.Join(r.GetClassTermPath(className), studentsFilename)
	data, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal students to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write students to file '%s': %w", filePath, err)
	}
	return nil
}

// LoadStudents loads students for a class/term
func (r *StudentRepository) LoadStudents(ctx context.Context, className string) ([]models.StudentInput, error) {
	if !r.IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}

	filePath := filepath.Join(r.GetClassTermPath(className), studentsFilename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("students file '%s' not found: %w", filePath, err)
		}
		return nil, fmt.Errorf("failed to read students file '%s': %w", filePath, err)
	}

	var students []models.StudentInput
	if err := json.Unmarshal(data, &students); err != nil {
		return nil, fmt.Errorf("failed to unmarshal students from JSON file '%s': %w", filePath, err)
	}
	return students, nil
}

// ReadStudentsFromFile reads students from a file
func (r *StudentRepository) ReadStudentsFromFile(ctx context.Context, fileName string) ([]models.Student, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var students []models.Student
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) == 0 {
			continue // Skip empty lines
		}

		name := strings.TrimSpace(parts[0])
		name = strings.Trim(name, "\"'")

		var preferences []string
		if len(parts) > 1 {
			for _, pref := range parts[1:] {
				trimmedPref := strings.TrimSpace(pref)
				trimmedPref = strings.Trim(trimmedPref, "\"'")
				preferences = append(preferences, trimmedPref)
			}
		}

		students = append(students, models.Student{
			Name:            name,
			Preferences:     preferences,
			Assigned:        false,
			AssignedProject: "",
			BadPreferences:  nil, // Initialize as nil, can be populated later
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
