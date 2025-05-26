package storage

import (
	"bufio"
	"os"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/models"
)

// readProjectsFromFile reads project names from a file.
// Each line in the file is expected to be a project name.
// Whitespace and surrounding quotes (single or double) are trimmed from each line.
func ReadProjectsFromFile(fileName string) ([]string, error) {
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

// readStudentsFromFile reads student data from a CSV file.
// Each line is expected to be in the format: Name,Preference1,Preference2,...
// Whitespace and surrounding quotes (single or double) are trimmed from each field.
func ReadStudentsFromFile(fileName string) ([]models.Student, error) {
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
