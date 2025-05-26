package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

const (
	projectsFilename = "projects.json"
)

// ProjectRepository implements the repository.ProjectRepository interface for file-based storage
type ProjectRepository struct {
	*BaseRepository
}

// NewProjectRepository creates a new file-based project repository
func NewProjectRepository(baseDir string) repository.ProjectRepository {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository(baseDir),
	}
}

// SaveProjects saves projects for a class/term
func (r *ProjectRepository) SaveProjects(ctx context.Context, className string, projects []string) error {
	if !r.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}

	if err := r.EnsureClassTermDir(className); err != nil {
		return err
	}

	filePath := filepath.Join(r.GetClassTermPath(className), projectsFilename)
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects to file '%s': %w", filePath, err)
	}
	return nil
}

// LoadProjects loads projects for a class/term
func (r *ProjectRepository) LoadProjects(ctx context.Context, className string) ([]string, error) {
	if !r.IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}

	filePath := filepath.Join(r.GetClassTermPath(className), projectsFilename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("projects file '%s' not found: %w", filePath, err)
		}
		return nil, fmt.Errorf("failed to read projects file '%s': %w", filePath, err)
	}

	var projects []string
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects from JSON file '%s': %w", filePath, err)
	}
	return projects, nil
}

// ReadProjectsFromFile reads projects from a file
func (r *ProjectRepository) ReadProjectsFromFile(ctx context.Context, fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var projects []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		// Remove surrounding quotes
		trimmedLine = strings.Trim(trimmedLine, "\"'")
		projects = append(projects, trimmedLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
