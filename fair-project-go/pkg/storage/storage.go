package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/models"
)

var (
	baseDataDir         = "data"
)

const (
	projectsFilename    = "projects.json"
	studentsFilename    = "students.json"
	assignmentFilePrefix = "assignment_"
)

// GetBaseDataDir returns the base data directory.
func GetBaseDataDir() string {
	return baseDataDir
}

// SetBaseDataDir sets the base data directory.
// This is useful for testing.
func SetBaseDataDir(dir string) {
	baseDataDir = dir
}

var classNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// IsValidClassName checks if the className is valid for use as a directory name.
// Allows alphanumeric characters, underscores, and hyphens.
func IsValidClassName(className string) bool {
	if className == "" {
		return false
	}
	return classNameRegex.MatchString(className)
}

// getClassTermPath constructs the path to a specific class/term directory.
func getClassTermPath(className string) string {
	return filepath.Join(baseDataDir, className)
}

// EnsureBaseDir checks if the base data directory (e.g., `data/`) exists.
// If not, it creates it.
func EnsureBaseDir() error {
	if _, err := os.Stat(baseDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDataDir, 0755); err != nil {
			return fmt.Errorf("failed to create base data directory '%s': %w", baseDataDir, err)
		}
		fmt.Printf("Base data directory '%s' created.\n", baseDataDir)
	} else if err != nil {
		return fmt.Errorf("failed to check base data directory '%s': %w", baseDataDir, err)
	}
	return nil
}

// CreateClassTermDir creates a new directory named className inside the base data directory.
func CreateClassTermDir(className string) error {
	if !IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	classTermPath := getClassTermPath(className)

	if err := EnsureBaseDir(); err != nil { // Ensure base directory exists first
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

// ListClassTerms reads the base data directory and returns a slice of strings
// containing the names of all class/term directories.
func ListClassTerms() ([]string, error) {
	if err := EnsureBaseDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(baseDataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read base data directory '%s': %w", baseDataDir, err)
	}

	var classTerms []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Further validation can be added here if non-class directories might exist
			if IsValidClassName(entry.Name()) { // Only list valid ones
				classTerms = append(classTerms, entry.Name())
			}
		}
	}
	return classTerms, nil
}

// SaveProjects saves the projects slice to data/className/projects.json.
func SaveProjects(className string, projects []string) error {
	if !IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}
	classTermPath := getClassTermPath(className)
	if _, err := os.Stat(classTermPath); os.IsNotExist(err) {
		return fmt.Errorf("class/term directory '%s' does not exist", classTermPath)
	}

	filePath := filepath.Join(classTermPath, projectsFilename)
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects to file '%s': %w", filePath, err)
	}
	return nil
}

// LoadProjects loads projects from data/className/projects.json.
func LoadProjects(className string) ([]string, error) {
	if !IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}
	filePath := filepath.Join(getClassTermPath(className), projectsFilename)

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

// SaveStudents saves the students slice to data/className/students.json.
func SaveStudents(className string, students []models.StudentInput) error {
	if !IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}
	classTermPath := getClassTermPath(className)
	if _, err := os.Stat(classTermPath); os.IsNotExist(err) {
		return fmt.Errorf("class/term directory '%s' does not exist", classTermPath)
	}

	filePath := filepath.Join(classTermPath, studentsFilename)
	data, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal students to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write students to file '%s': %w", filePath, err)
	}
	return nil
}

// LoadStudents loads students from data/className/students.json.
func LoadStudents(className string) ([]models.StudentInput, error) {
	if !IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}
	filePath := filepath.Join(getClassTermPath(className), studentsFilename)

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

// SaveAssignmentResults saves the assignment results to data/className/assignment_<assignmentID>.json.
func SaveAssignmentResults(className string, assignmentID string, results models.AssignmentOutput) error {
	if !IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s'", className)
	}
	if !IsValidClassName(assignmentID) { // Reuse classNameRegex for assignmentID for simplicity
		return fmt.Errorf("invalid assignment ID '%s': must be alphanumeric, underscore, or hyphen", assignmentID)
	}

	classTermPath := getClassTermPath(className)
	if _, err := os.Stat(classTermPath); os.IsNotExist(err) {
		return fmt.Errorf("class/term directory '%s' does not exist", classTermPath)
	}

	fileName := fmt.Sprintf("%s%s.json", assignmentFilePrefix, assignmentID)
	filePath := filepath.Join(classTermPath, fileName)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal assignment results to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write assignment results to file '%s': %w", filePath, err)
	}
	return nil
}

// ListAssignmentResults lists all assignment_*.json files in the data/className/ directory.
// Returns a slice of assignment IDs.
func ListAssignmentResults(className string) ([]string, error) {
	if !IsValidClassName(className) {
		return nil, fmt.Errorf("invalid class/term name '%s'", className)
	}
	classTermPath := getClassTermPath(className)
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
			if IsValidClassName(id) { // Validate the extracted ID part
				assignmentIDs = append(assignmentIDs, id)
			}
		}
	}
	return assignmentIDs, nil
}

// LoadAssignmentResult loads a specific assignment result from data/className/assignment_<assignmentID>.json.
func LoadAssignmentResult(className string, assignmentID string) (models.AssignmentOutput, error) {
	var results models.AssignmentOutput
	if !IsValidClassName(className) {
		return results, fmt.Errorf("invalid class/term name '%s'", className)
	}
	if !IsValidClassName(assignmentID) {
		return results, fmt.Errorf("invalid assignment ID '%s'", assignmentID)
	}

	fileName := fmt.Sprintf("%s%s.json", assignmentFilePrefix, assignmentID)
	filePath := filepath.Join(getClassTermPath(className), fileName)

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
