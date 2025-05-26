package file

import (
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
	assignmentFilePrefix = "assignment_"
)

// AssignmentRepository implements the repository.AssignmentRepository interface for file-based storage
type AssignmentRepository struct {
	*BaseRepository
}

// NewAssignmentRepository creates a new file-based assignment repository
func NewAssignmentRepository(baseDir string) repository.AssignmentRepository {
	return &AssignmentRepository{
		BaseRepository: NewBaseRepository(baseDir),
	}
}

// SaveAssignmentResults saves assignment results for a class/term
func (r *AssignmentRepository) SaveAssignmentResults(ctx context.Context, className string, assignmentID string, results models.AssignmentOutput) error {
	if !r.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}
	if !r.IsValidClassName(assignmentID) { // Reuse classNameRegex for assignmentID for simplicity
		return fmt.Errorf("invalid assignment ID '%s': must be alphanumeric, underscore, or hyphen", assignmentID)
	}

	if err := r.EnsureClassTermDir(className); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s%s.json", assignmentFilePrefix, assignmentID)
	filePath := filepath.Join(r.GetClassTermPath(className), fileName)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal assignment results to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write assignment results to file '%s': %w", filePath, err)
	}
	return nil
}

// ListAssignmentResults lists all assignment results for a class/term
func (r *AssignmentRepository) ListAssignmentResults(ctx context.Context, className string) ([]string, error) {
	if !r.IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}

	classTermPath := r.GetClassTermPath(className)
	if _, err := os.Stat(classTermPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("class/term directory '%s' does not exist", classTermPath)
	}

	entries, err := os.ReadDir(classTermPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read class/term directory '%s': %w", classTermPath, err)
	}

	var assignmentIDs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), assignmentFilePrefix) && strings.HasSuffix(entry.Name(), ".json") {
			id := strings.TrimPrefix(entry.Name(), assignmentFilePrefix)
			id = strings.TrimSuffix(id, ".json")
			if r.IsValidClassName(id) { // Validate the extracted ID part
				assignmentIDs = append(assignmentIDs, id)
			}
		}
	}
	return assignmentIDs, nil
}

// LoadAssignmentResult loads an assignment result for a class/term
func (r *AssignmentRepository) LoadAssignmentResult(ctx context.Context, className string, assignmentID string) (models.AssignmentOutput, error) {
	var results models.AssignmentOutput
	if !r.IsValidClassName(className) {
		return results, fmt.Errorf("invalid class/term name '%s'", className)
	}
	if !r.IsValidClassName(assignmentID) {
		return results, fmt.Errorf("invalid assignment ID '%s'", assignmentID)
	}

	fileName := fmt.Sprintf("%s%s.json", assignmentFilePrefix, assignmentID)
	filePath := filepath.Join(r.GetClassTermPath(className), fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return results, fmt.Errorf("assignment results file '%s' not found: %w", filePath, err)
		}
		return results, fmt.Errorf("failed to read assignment results file '%s': %w", filePath, err)
	}

	if err := json.Unmarshal(data, &results); err != nil {
		return results, fmt.Errorf("failed to unmarshal assignment results from JSON file '%s': %w", filePath, err)
	}
	return results, nil
}
