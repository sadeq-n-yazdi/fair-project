package file

import (
	"context"
	"fmt"
	"os"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// ClassRepository implements the repository.ClassRepository interface for file-based storage
type ClassRepository struct {
	*BaseRepository
}

// NewClassRepository creates a new file-based class repository
func NewClassRepository(baseDir string) repository.ClassRepository {
	return &ClassRepository{
		BaseRepository: NewBaseRepository(baseDir),
	}
}

// CreateClassTerm creates a new class/term directory
func (r *ClassRepository) CreateClassTerm(ctx context.Context, className string) error {
	if !r.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	classTermPath := r.GetClassTermPath(className)

	if err := r.EnsureBaseDir(); err != nil { // Ensure base directory exists first
		return err
	}

	// Attempt to create the directory. os.Mkdir will fail if it already exists.
	// And its error can be checked with os.IsExist.
	err := os.Mkdir(classTermPath, 0755)
	if err != nil {
		// Wrap the error to allow os.IsExist and other checks in handlers
		return fmt.Errorf("failed to create class/term directory '%s': %w", classTermPath, err)
	}
	fmt.Printf("Class/term directory '%s' created.\n", classTermPath)
	return nil
}

// ListClassTerms lists all class/terms
func (r *ClassRepository) ListClassTerms(ctx context.Context) ([]string, error) {
	if err := r.EnsureBaseDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read base data directory '%s': %w", r.baseDir, err)
	}

	var classTerms []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Further validation can be added here if non-class directories might exist
			if r.IsValidClassName(entry.Name()) { // Only list valid ones
				classTerms = append(classTerms, entry.Name())
			}
		}
	}
	return classTerms, nil
}
