package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// createClassTermHandler handles POST requests to /classes
// Request: JSON body {"name": "className"}
// Response: Success (201) or Error (400, 409, 500)
func createClassTermHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var reqBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	className := strings.TrimSpace(reqBody.Name)
	if className == "" {
		respondError(w, http.StatusBadRequest, "Class/term name cannot be empty")
		return
	}

	// isValidClassName is in storage.go, assumed accessible.
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format. Use alphanumeric, underscores, or hyphens.")
		return
	}

	err := CreateClassTermDir(className) // From storage.go
	if err != nil {
		if os.IsExist(err) { // This check should now work due to error wrapping in CreateClassTermDir
			respondError(w, http.StatusConflict, "Class/term directory '"+className+"' already exists.")
		} else if strings.Contains(err.Error(), "invalid class/term name") {
			// Check for the specific validation error text from CreateClassTermDir
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to create class/term directory: "+err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Class/term '" + className + "' created successfully"})
}

// listClassTermsHandler handles GET requests to /classes
// Response: Success (200) with list of class names or Error (500)
func listClassTermsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	classTerms, err := ListClassTerms() // From storage.go
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list class/terms: "+err.Error())
		return
	}

	if classTerms == nil { // Ensure we return an empty list, not null, if no terms exist
		classTerms = []string{}
	}
	respondJSON(w, http.StatusOK, classTerms)
}

// projectsHandler manages project-related requests for a specific class/term.
// It routes based on method (GET for loading, POST for saving).
// URL: /classes/{className}/projects
func projectsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL.")
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload projects
		var projects []string
		if err := json.NewDecoder(r.Body).Decode(&projects); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body, expected JSON array of project strings: "+err.Error())
			return
		}
		defer r.Body.Close()

		// Basic validation: check if projects list is empty (optional, based on requirements)
		// if len(projects) == 0 {
		// respondError(w, http.StatusBadRequest, "Projects list cannot be empty")
		// return
		// }

		err := SaveProjects(className, projects) // From storage.go
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				respondError(w, http.StatusNotFound, "Class/term directory '"+className+"' not found: "+err.Error())
			} else {
				respondError(w, http.StatusInternalServerError, "Failed to save projects: "+err.Error())
			}
			return
		}
		respondJSON(w, http.StatusCreated, map[string]string{"message": "Projects for '" + className + "' saved successfully"})

	case http.MethodGet: // Get projects
		projects, err := LoadProjects(className) // From storage.go
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, http.StatusNotFound, "Projects file not found for class/term '"+className+"': "+err.Error())
			} else {
				respondError(w, http.StatusInternalServerError, "Failed to load projects: "+err.Error())
			}
			return
		}
		if projects == nil { // Ensure empty list, not null
			projects = []string{}
		}
		respondJSON(w, http.StatusOK, projects)

	default:
		respondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for projects")
	}
}

// studentsHandler manages student-related requests for a specific class/term.
// It routes based on method (GET for loading, POST for saving).
// URL: /classes/{className}/students
func studentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL.")
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload students
		var students []StudentInput // Type from types.go
		if err := json.NewDecoder(r.Body).Decode(&students); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body, expected JSON array of student objects: "+err.Error())
			return
		}
		defer r.Body.Close()

		// Add any validation for student data here if needed (e.g., non-empty list, valid student names)

		err := SaveStudents(className, students) // From storage.go
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				respondError(w, http.StatusNotFound, "Class/term directory '"+className+"' not found: "+err.Error())
			} else {
				respondError(w, http.StatusInternalServerError, "Failed to save students: "+err.Error())
			}
			return
		}
		respondJSON(w, http.StatusCreated, map[string]string{"message": "Students for '" + className + "' saved successfully"})

	case http.MethodGet: // Get students
		students, err := LoadStudents(className) // From storage.go
		if err != nil {
			if os.IsNotExist(err) {
				respondError(w, http.StatusNotFound, "Students file not found for class/term '"+className+"': "+err.Error())
			} else {
				respondError(w, http.StatusInternalServerError, "Failed to load students: "+err.Error())
			}
			return
		}
		if students == nil { // Ensure empty list, not null
			students = []StudentInput{}
		}
		respondJSON(w, http.StatusOK, students)

	default:
		respondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for students")
	}
}

// masterClassResourceHandler is a generic handler for /classes/{className}/{resource}
// e.g., /classes/Fall2023/projects or /classes/Spring2024/students
// It will delegate to specific handlers like projectsHandler or studentsHandler.
func masterClassResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/classes/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		// This case should ideally be handled by a more specific /classes handler if it's for listing or creating classes.
		// If it reaches here, it means it's an invalid path like /classes//projects
		respondError(w, http.StatusBadRequest, "Class/term name missing in URL path")
		return
	}
	className := parts[0]

	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL: "+className)
		return
	}

	if len(parts) == 1 {
		// Potentially a handler for /classes/{className} - e.g., get class details (not specified yet)
		respondError(w, http.StatusNotFound, "Resource type (e.g., projects, students) missing in URL path for class '"+className+"'")
		return
	}

	resourceType := parts[1]
	switch resourceType {
	case "projects":
		projectsHandler(w, r, className)
	case "students":
		studentsHandler(w, r, className)
	case "assign":
		// triggerAssignmentHandler is defined in assignment_handlers.go
		// It expects a POST request.
		triggerAssignmentHandler(w, r, className) // This only handles POST
	case "assignments":
		// Path: /classes/{className}/assignments OR /classes/{className}/assignments/{assignmentID}
		if len(parts) == 2 { // Exactly /classes/{className}/assignments
			listAssignmentsHandler(w, r, className) // Defined in assignment_handlers.go
		} else if len(parts) == 3 { // Exactly /classes/{className}/assignments/{assignmentID}
			assignmentID := parts[2]
			if !isValidClassName(assignmentID) { // Reuse validation for simplicity
				respondError(w, http.StatusBadRequest, "Invalid assignment ID format in URL: "+assignmentID)
				return
			}
			getAssignmentResultHandler(w, r, className, assignmentID) // Defined in assignment_handlers.go
		} else {
			// Path is too long, e.g., /classes/{className}/assignments/{id}/something_else
			respondError(w, http.StatusNotFound, "Invalid path structure for assignments resource.")
		}
	default:
		respondError(w, http.StatusNotFound, "Unknown resource type '"+resourceType+"' for class '"+className+"'")
	}
}
