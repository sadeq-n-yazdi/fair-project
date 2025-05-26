package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// BaseRepository provides common functionality for file-based repositories
type BaseRepository struct {
	baseDir string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(baseDir string) *BaseRepository {
	return &BaseRepository{
		baseDir: baseDir,
	}
}

// Initialize initializes the repository
func (r *BaseRepository) Initialize(_ context.Context) error {
	return r.EnsureBaseDir()
}

// Close closes the repository
func (r *BaseRepository) Close(_ context.Context) error {
	// No resources to close for file-based storage
	return nil
}

// EnsureBaseDir checks if the base data directory exists and creates it if it doesn't
func (r *BaseRepository) EnsureBaseDir() error {
	if _, err := os.Stat(r.baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(r.baseDir, 0755); err != nil {
			return fmt.Errorf("failed to create base data directory '%s': %w", r.baseDir, err)
		}
		fmt.Printf("Base data directory '%s' created.\n", r.baseDir)
	} else if err != nil {
		return fmt.Errorf("failed to check base data directory '%s': %w", r.baseDir, err)
	}
	return nil
}

var classNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// IsValidClassName checks if the className is valid for use as a directory name
func (r *BaseRepository) IsValidClassName(className string) bool {
	if className == "" {
		return false
	}
	return classNameRegex.MatchString(className)
}

// GetClassTermPath constructs the path to a specific class/term directory
func (r *BaseRepository) GetClassTermPath(className string) string {
	return filepath.Join(r.baseDir, className)
}

// EnsureClassTermDir ensures that the class/term directory exists
func (r *BaseRepository) EnsureClassTermDir(className string) error {
	if !r.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	classTermPath := r.GetClassTermPath(className)
	if _, err := os.Stat(classTermPath); os.IsNotExist(err) {
		return fmt.Errorf("class/term directory '%s' does not exist", classTermPath)
	}
	return nil
}
