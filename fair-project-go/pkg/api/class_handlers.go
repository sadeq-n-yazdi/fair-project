package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// CreateClassTermHandler handles POST requests to /classes
// Request: JSON body {"name": "className"}
// Response: Success (201) or Error (400, 409, 500)
func CreateClassTermHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var reqBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	className := strings.TrimSpace(reqBody.Name)
	if className == "" {
		RespondError(w, http.StatusBadRequest, "Class/term name cannot be empty")
		return
	}

	// IsValidClassName is in storage.go
	if !storage.IsValidClassName(className) {
		RespondError(w, http.StatusBadRequest, "Invalid class/term name format. Use alphanumeric, underscores, or hyphens.")
		return
	}

	err := storage.CreateClassTermDir(className) // From storage.go
	if err != nil {
		if os.IsExist(err) { // This check should now work due to error wrapping in CreateClassTermDir
			RespondError(w, http.StatusConflict, "Class/term directory '"+className+"' already exists.")
		} else if strings.Contains(err.Error(), "invalid class/term name") {
			// Check for the specific validation error text from CreateClassTermDir
			RespondError(w, http.StatusBadRequest, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to create class/term directory: "+err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]string{"message": "Class/term '" + className + "' created successfully"})
}

// ListClassTermsHandler handles GET requests to /classes
// Response: Success (200) with list of class names or Error (500)
func ListClassTermsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	classTerms, err := storage.ListClassTerms() // From storage.go
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list class/terms: "+err.Error())
		return
	}

	if classTerms == nil { // Ensure we return an empty list, not null, if no terms exist
		classTerms = []string{}
	}
	RespondJSON(w, http.StatusOK, classTerms)
}

// ProjectsHandler manages project-related requests for a specific class/term.
// It routes based on method (GET for loading, POST for saving).
// URL: /classes/{className}/projects
func ProjectsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !storage.IsValidClassName(className) {
		RespondError(w, http.StatusBadRequest, "Invalid class/term name format in URL.")
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload projects
		var projects []string
		if err := json.NewDecoder(r.Body).Decode(&projects); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid request body, expected JSON array of project strings: "+err.Error())
			return
		}
		defer r.Body.Close()

		// Basic validation: check if projects list is empty (optional, based on requirements)
		// if len(projects) == 0 {
		// RespondError(w, http.StatusBadRequest, "Projects list cannot be empty")
		// return
		// }

		err := storage.SaveProjects(className, projects) // From storage.go
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				RespondError(w, http.StatusNotFound, "Class/term directory '"+className+"' not found: "+err.Error())
			} else {
				RespondError(w, http.StatusInternalServerError, "Failed to save projects: "+err.Error())
			}
			return
		}
		RespondJSON(w, http.StatusCreated, map[string]string{"message": "Projects for '" + className + "' saved successfully"})

	case http.MethodGet: // Get projects
		projects, err := storage.LoadProjects(className) // From storage.go
		if err != nil {
			if os.IsNotExist(err) {
				RespondError(w, http.StatusNotFound, "Projects file not found for class/term '"+className+"': "+err.Error())
			} else {
				RespondError(w, http.StatusInternalServerError, "Failed to load projects: "+err.Error())
			}
			return
		}
		if projects == nil { // Ensure empty list, not null
			projects = []string{}
		}
		RespondJSON(w, http.StatusOK, projects)

	default:
		RespondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for projects")
	}
}

// StudentsHandler manages student-related requests for a specific class/term.
// It routes based on method (GET for loading, POST for saving).
// URL: /classes/{className}/students
func StudentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !storage.IsValidClassName(className) {
		RespondError(w, http.StatusBadRequest, "Invalid class/term name format in URL.")
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload students
		var students []models.StudentInput // Type from models package
		if err := json.NewDecoder(r.Body).Decode(&students); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid request body, expected JSON array of student objects: "+err.Error())
			return
		}
		defer r.Body.Close()

		// Add any validation for student data here if needed (e.g., non-empty list, valid student names)

		err := storage.SaveStudents(className, students) // From storage.go
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				RespondError(w, http.StatusNotFound, "Class/term directory '"+className+"' not found: "+err.Error())
			} else {
				RespondError(w, http.StatusInternalServerError, "Failed to save students: "+err.Error())
			}
			return
		}
		RespondJSON(w, http.StatusCreated, map[string]string{"message": "Students for '" + className + "' saved successfully"})

	case http.MethodGet: // Get students
		students, err := storage.LoadStudents(className) // From storage.go
		if err != nil {
			if os.IsNotExist(err) {
				RespondError(w, http.StatusNotFound, "Students file not found for class/term '"+className+"': "+err.Error())
			} else {
				RespondError(w, http.StatusInternalServerError, "Failed to load students: "+err.Error())
			}
			return
		}
		if students == nil { // Ensure empty list, not null
			students = []models.StudentInput{}
		}
		RespondJSON(w, http.StatusOK, students)

	default:
		RespondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for students")
	}
}

// MasterClassResourceHandler is a generic handler for /classes/{className}/{resource}
// e.g., /classes/Fall2023/projects or /classes/Spring2024/students
// It will delegate to specific handlers like ProjectsHandler or StudentsHandler.
func MasterClassResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/classes/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		// This case should ideally be handled by a more specific /classes handler if it's for listing or creating classes.
		// If it reaches here, it means it's an invalid path like /classes//projects
		RespondError(w, http.StatusBadRequest, "Class/term name missing in URL path")
		return
	}
	className := parts[0]

	if !storage.IsValidClassName(className) {
		RespondError(w, http.StatusBadRequest, "Invalid class/term name format in URL: "+className)
		return
	}

	if len(parts) == 1 {
		// Potentially a handler for /classes/{className} - e.g., get class details (not specified yet)
		RespondError(w, http.StatusNotFound, "Resource type (e.g., projects, students) missing in URL path for class '"+className+"'")
		return
	}

	resourceType := parts[1]
	switch resourceType {
	case "projects":
		ProjectsHandler(w, r, className)
	case "students":
		StudentsHandler(w, r, className)
	case "assign":
		// TriggerAssignmentHandler is defined in assignment_handlers.go
		// It expects a POST request.
		TriggerAssignmentHandler(w, r, className) // This only handles POST
	case "assignments":
		// Path: /classes/{className}/assignments OR /classes/{className}/assignments/{assignmentID}
		if len(parts) == 2 { // Exactly /classes/{className}/assignments
			ListAssignmentsHandler(w, r, className) // Defined in assignment_handlers.go
		} else if len(parts) == 3 { // Exactly /classes/{className}/assignments/{assignmentID}
			assignmentID := parts[2]
			if !storage.IsValidClassName(assignmentID) { // Reuse validation for simplicity
				RespondError(w, http.StatusBadRequest, "Invalid assignment ID format in URL: "+assignmentID)
				return
			}
			GetAssignmentResultHandler(w, r, className, assignmentID) // Defined in assignment_handlers.go
		} else {
			// Path is too long, e.g., /classes/{className}/assignments/{id}/something_else
			RespondError(w, http.StatusNotFound, "Invalid path structure for assignments resource.")
		}
	default:
		RespondError(w, http.StatusNotFound, "Unknown resource type '"+resourceType+"' for class '"+className+"'")
	}
}
