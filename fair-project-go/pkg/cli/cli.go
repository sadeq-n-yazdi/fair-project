package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sadeq/fair-project-go/pkg/assignment"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// CreateClassTerm creates a new class/term directory
func CreateClassTerm(className string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	err := storage.CreateClassTermDir(className)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("class/term directory '%s' already exists", className)
		}
		return fmt.Errorf("failed to create class/term directory: %v", err)
	}

	fmt.Printf("Class/term '%s' created successfully\n", className)
	return nil
}

// ListClassTerms lists all class/term directories
func ListClassTerms() error {
	classTerms, err := storage.ListClassTerms()
	if err != nil {
		return fmt.Errorf("failed to list class/terms: %v", err)
	}

	if len(classTerms) == 0 {
		fmt.Println("No class/terms found")
		return nil
	}

	fmt.Println("Available class/terms:")
	for _, className := range classTerms {
		fmt.Println("-", className)
	}
	return nil
}

// ImportProjects imports projects from a file
func ImportProjects(className string, filePath string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	projects, err := storage.ReadProjectsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read projects from file: %v", err)
	}

	err = storage.SaveProjects(className, projects)
	if err != nil {
		return fmt.Errorf("failed to save projects: %v", err)
	}

	fmt.Printf("Imported %d projects for class/term '%s'\n", len(projects), className)
	return nil
}

// ImportStudents imports students from a file
func ImportStudents(className string, filePath string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	students, err := storage.ReadStudentsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read students from file: %v", err)
	}

	// Convert []models.Student to []models.StudentInput
	var studentInputs []models.StudentInput
	for _, student := range students {
		studentInputs = append(studentInputs, models.StudentInput{
			Name:        student.Name,
			Preferences: student.Preferences,
		})
	}

	err = storage.SaveStudents(className, studentInputs)
	if err != nil {
		return fmt.Errorf("failed to save students: %v", err)
	}

	fmt.Printf("Imported %d students for class/term '%s'\n", len(studentInputs), className)
	return nil
}

// ListProjects lists all projects for a class/term
func ListProjects(className string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	projects, err := storage.LoadProjects(className)
	if err != nil {
		return fmt.Errorf("failed to load projects: %v", err)
	}

	if len(projects) == 0 {
		fmt.Printf("No projects found for class/term '%s'\n", className)
		return nil
	}

	fmt.Printf("Projects for class/term '%s':\n", className)
	for i, project := range projects {
		fmt.Printf("%d. %s\n", i+1, project)
	}
	return nil
}

// ListStudents lists all students for a class/term
func ListStudents(className string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	students, err := storage.LoadStudents(className)
	if err != nil {
		return fmt.Errorf("failed to load students: %v", err)
	}

	if len(students) == 0 {
		fmt.Printf("No students found for class/term '%s'\n", className)
		return nil
	}

	fmt.Printf("Students for class/term '%s':\n", className)
	for i, student := range students {
		fmt.Printf("%d. %s (Preferences: %v)\n", i+1, student.Name, student.Preferences)
	}
	return nil
}

// RunAssignment runs the assignment algorithm for a class/term
func RunAssignment(className string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	// Load projects
	projects, err := storage.LoadProjects(className)
	if err != nil {
		return fmt.Errorf("failed to load projects: %v", err)
	}
	if len(projects) == 0 {
		return fmt.Errorf("no projects found for class/term '%s'", className)
	}

	// Load students
	studentInputs, err := storage.LoadStudents(className)
	if err != nil {
		return fmt.Errorf("failed to load students: %v", err)
	}
	if len(studentInputs) == 0 {
		return fmt.Errorf("no students found for class/term '%s'", className)
	}

	// Convert []StudentInput to []Student
	var students []models.Student
	for _, si := range studentInputs {
		students = append(students, models.Student{
			Name:            si.Name,
			Preferences:     si.Preferences,
			Assigned:        false,
			AssignedProject: "",
			BadPreferences:  []string{},
		})
	}

	// Run assignment algorithm
	selectedByProjects, selectedByStudent, updatedStudents, remainingProjects := assignment.AssignProjects(students, projects)

	// Create assignment output
	assignmentOutput := models.AssignmentOutput{
		SelectedByProjects: selectedByProjects,
		SelectedByStudent:  selectedByStudent,
		Students:           updatedStudents,
		RemainingProjects:  remainingProjects,
	}

	// Generate assignment ID
	assignmentID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Save assignment results
	err = storage.SaveAssignmentResults(className, assignmentID, assignmentOutput)
	if err != nil {
		return fmt.Errorf("failed to save assignment results: %v", err)
	}

	fmt.Printf("Assignment completed successfully for class/term '%s'\n", className)
	fmt.Printf("Assignment ID: %s\n", assignmentID)
	fmt.Printf("Assigned %d students to projects\n", len(selectedByStudent))
	fmt.Printf("Remaining projects: %d\n", len(remainingProjects))

	return nil
}

// ListAssignments lists all assignments for a class/term
func ListAssignments(className string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}

	assignmentIDs, err := storage.ListAssignmentResults(className)
	if err != nil {
		return fmt.Errorf("failed to list assignments: %v", err)
	}

	if len(assignmentIDs) == 0 {
		fmt.Printf("No assignments found for class/term '%s'\n", className)
		return nil
	}

	fmt.Printf("Assignments for class/term '%s':\n", className)
	for i, assignmentID := range assignmentIDs {
		fmt.Printf("%d. %s\n", i+1, assignmentID)
	}
	return nil
}

// ShowAssignment shows the details of an assignment
func ShowAssignment(className string, assignmentID string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}
	if !storage.IsValidClassName(assignmentID) {
		return fmt.Errorf("invalid assignment ID '%s': must be alphanumeric, underscore, or hyphen", assignmentID)
	}

	result, err := storage.LoadAssignmentResult(className, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to load assignment result: %v", err)
	}

	fmt.Printf("Assignment '%s' for class/term '%s':\n", assignmentID, className)
	fmt.Println("Assigned students:")
	for student, project := range result.SelectedByStudent {
		fmt.Printf("- %s -> %s\n", student, project)
	}

	fmt.Println("\nRemaining projects:")
	for _, project := range result.RemainingProjects {
		fmt.Printf("- %s\n", project)
	}

	return nil
}

// ExportAssignment exports an assignment to a JSON file
func ExportAssignment(className string, assignmentID string, outputPath string) error {
	if !storage.IsValidClassName(className) {
		return fmt.Errorf("invalid class/term name '%s': must be alphanumeric, underscore, or hyphen", className)
	}
	if !storage.IsValidClassName(assignmentID) {
		return fmt.Errorf("invalid assignment ID '%s': must be alphanumeric, underscore, or hyphen", assignmentID)
	}

	result, err := storage.LoadAssignmentResult(className, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to load assignment result: %v", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Marshal assignment result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal assignment result to JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write assignment result to file: %v", err)
	}

	fmt.Printf("Assignment '%s' for class/term '%s' exported to %s\n", assignmentID, className, outputPath)
	return nil
}
