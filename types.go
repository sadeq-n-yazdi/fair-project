package main

// The data for a specific class or term will be stored in a directory
// named after that class/term, located within a base 'data/' directory.
// For example, data for "Fall2024ComputerScience" would be in
// "data/Fall2024ComputerScience/".
//
// Class/term names should be simple, URL-safe strings (e.g., alphanumeric,
// underscores, hyphens) to avoid issues with file paths and potential API routes.

// Project is a type alias for string, representing the name of a project.
// Input file: projects.json - A JSON array of strings.
// Example: ["Project_Alpha", "Project_Beta"]
type Project string

// StudentInput is used for unmarshalling student data from students.json.
// Input file: students.json - A JSON array of StudentInput objects.
type StudentInput struct {
	Name        string   `json:"name"`
	Preferences []string `json:"preferences"`
}

// Student represents a student with their preferences and assignment status.
// This struct is used internally and for outputting assignment results.
type Student struct {
	Name            string   `json:"Name"` // Capitalized for consistency with potential existing uses if any, or can be `name`
	Preferences     []string `json:"Preferences"`
	Assigned        bool     `json:"Assigned"`
	AssignedProject string   `json:"AssignedProject"`
	BadPreferences  []string `json:"BadPreferences"`
}

// AssignmentOutput represents the structure for storing the results of an assignment run.
// Output file: assignment_<ID>.json - A JSON object.
type AssignmentOutput struct {
	SelectedByProjects map[string]string `json:"selectedByProjects"`
	SelectedByStudent  map[string]string `json:"selectedByStudent"`
	Students           []Student         `json:"students"` // Note: uses the main Student struct
	RemainingProjects  []string          `json:"remainingProjects"`
}
