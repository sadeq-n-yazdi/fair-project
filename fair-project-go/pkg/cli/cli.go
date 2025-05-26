package cli

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sadeq/fair-project-go/pkg/assignment"
	"github.com/sadeq/fair-project-go/pkg/auth"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// generateRandomString generates a random alphanumeric string of the specified length
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// updateEnvFile updates or adds a key-value pair in the .env file
func updateEnvFile(key, value string) error {
	// Try to open the .env file
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try to find the .env file in the parent directory
		envPath = "../.env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			return fmt.Errorf("could not find .env file")
		}
	}

	// Read the file content
	content, err := ioutil.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %v", err)
	}

	// Create a regex to find the key
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^%s=.*$`, regexp.QuoteMeta(key)))

	// Check if the key exists in the file
	if re.Match(content) {
		// Replace the existing key-value pair
		newContent := re.ReplaceAllString(string(content), fmt.Sprintf("%s=%s", key, value))
		err = ioutil.WriteFile(envPath, []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write to .env file: %v", err)
		}
	} else {
		// Append the new key-value pair to the file
		file, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .env file: %v", err)
		}
		defer file.Close()

		// Add a newline if the file doesn't end with one
		if len(content) > 0 && content[len(content)-1] != '\n' {
			_, err = file.WriteString("\n")
			if err != nil {
				return fmt.Errorf("failed to write to .env file: %v", err)
			}
		}

		// Write the new key-value pair
		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("failed to write to .env file: %v", err)
		}
	}

	return nil
}

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

// CreateSuperAdmin creates the first superadmin user
// This function requires the SUPERADMIN_KEY from the .env file
func CreateSuperAdmin(username, password, key string) error {
	// Check if the superadmin key is correct
	expectedKey := os.Getenv("SUPERADMIN_KEY")
	if expectedKey == "" {
		return fmt.Errorf("SUPERADMIN_KEY not set in environment variables")
	}

	if key != expectedKey {
		return fmt.Errorf("invalid superadmin key")
	}

	// Check if any users already exist
	users, err := storage.LoadUsers()
	if err != nil {
		return fmt.Errorf("failed to load users: %v", err)
	}

	// Check if there are any superadmins already
	for _, user := range users {
		for _, role := range user.Roles {
			if role == models.RoleSuperAdmin {
				return fmt.Errorf("a superadmin user already exists")
			}
		}
	}

	// Hash the password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create the superadmin user
	superadmin := models.User{
		Username:     username,
		PasswordHash: passwordHash,
		Roles:        []models.Role{models.RoleSuperAdmin},
		Enabled:      true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the user
	if err := storage.SaveUser(superadmin); err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	fmt.Printf("Superadmin user '%s' created successfully\n", username)
	return nil
}

// ChangeUserPassword changes a user's password
func ChangeUserPassword(username, newPassword string) error {
	// Get the user
	user, exists, err := storage.GetUser(username)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if !exists {
		return fmt.Errorf("user '%s' not found", username)
	}

	// Hash the new password
	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update the user's password
	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()

	// Save the user
	if err := storage.SaveUser(user); err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	fmt.Printf("Password for user '%s' changed successfully\n", username)
	return nil
}

// ListUsers lists all users
func ListUsers() error {
	// Get all users
	users, err := storage.LoadUsers()
	if err != nil {
		return fmt.Errorf("failed to load users: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	fmt.Println("Users:")
	for _, user := range users {
		status := "Enabled"
		if !user.Enabled {
			status = "Disabled"
		}
		fmt.Printf("- %s (Roles: %v, Status: %s)\n", user.Username, user.Roles, status)
	}
	return nil
}

// PromptForSuperAdmin prompts the user to create a superadmin if none exists
func PromptForSuperAdmin() error {
	// Check if any users already exist
	users, err := storage.LoadUsers()
	if err != nil {
		return fmt.Errorf("failed to load users: %v", err)
	}

	// Check if there are any superadmins already
	for _, user := range users {
		for _, role := range user.Roles {
			if role == models.RoleSuperAdmin {
				// Superadmin already exists, no need to prompt
				return nil
			}
		}
	}

	// No superadmin exists, prompt to create one
	fmt.Println("No superadmin user found. Would you like to create one? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Skipping superadmin creation")
		return nil
	}

	// Get superadmin details
	fmt.Print("Enter username for superadmin: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Enter password for superadmin: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	password = strings.TrimSpace(password)

	// Generate a random superadmin key
	key, err := generateRandomString(32)
	if err != nil {
		return fmt.Errorf("failed to generate superadmin key: %v", err)
	}

	// Update the .env file with the generated key
	err = updateEnvFile("SUPERADMIN_KEY", key)
	if err != nil {
		return fmt.Errorf("failed to update .env file: %v", err)
	}

	// Show the generated key to the user
	fmt.Printf("Generated SUPERADMIN_KEY: %s\n", key)
	fmt.Println("This key has been saved to your .env file.")
	fmt.Println("Please keep it secure as it will be needed for future superadmin operations.")

	// Create the superadmin
	return CreateSuperAdmin(username, password, key)
}

// GenerateEnvFile generates a new .env file with random values for sensitive fields
func GenerateEnvFile(outputPath string) error {
	// Generate random values for sensitive fields
	jwtSecret, err := generateRandomString(64)
	if err != nil {
		return fmt.Errorf("failed to generate JWT secret: %v", err)
	}

	superadminKey, err := generateRandomString(32)
	if err != nil {
		return fmt.Errorf("failed to generate superadmin key: %v", err)
	}

	postgresPassword, err := generateRandomString(16)
	if err != nil {
		return fmt.Errorf("failed to generate PostgreSQL password: %v", err)
	}

	// Create the .env file content
	envContent := `# Fair Project API Configuration

# JWT Secret for authentication
JWT_SECRET=%s

# Super Admin Key (used for creating the first super admin via CLI)
SUPERADMIN_KEY=%s

# Server port
PORT=8080

# Data directory
DATA_DIR=data

# Log level (NONE, ERROR, INFO, DEBUG)
LOG_LEVEL=INFO

# PostgreSQL Configuration
POSTGRES_USER=fairuser
POSTGRES_PASSWORD=%s
POSTGRES_DB=fairdb
POSTGRES_PORT=5432
`

	// Format the content with the generated values
	envContent = fmt.Sprintf(envContent, jwtSecret, superadminKey, postgresPassword)

	// Write the content to the output file
	err = os.WriteFile(outputPath, []byte(envContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .env file: %v", err)
	}

	fmt.Printf("Generated .env file at %s with random values for sensitive fields.\n", outputPath)
	fmt.Println("Please keep this file secure as it contains sensitive information.")
	return nil
}
