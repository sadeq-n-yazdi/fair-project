package api

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/errors"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// Handler is the main API handler that holds a reference to the storage manager
type Handler struct {
	storageManager *storage.Manager
}

// NewHandler creates a new API handler with the given storage manager
func NewHandler(storageManager *storage.Manager) *Handler {
	return &Handler{
		storageManager: storageManager,
	}
}

// CreateClassTermHandler handles POST requests to /classes
// Request: JSON body {"name": "className"}
// Response: Success (201) or Error (400, 409, 500)
func (h *Handler) CreateClassTermHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only POST method is allowed"))
		return
	}

	var reqBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body"))
		return
	}
	defer r.Body.Close()

	className := strings.TrimSpace(reqBody.Name)
	if className == "" {
		HandleError(w, errors.New(errors.ErrEmptyClassName, "Class/term name cannot be empty"))
		return
	}

	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format. Use alphanumeric, underscores, or hyphens"))
		return
	}

	// Get the class repository from the storage manager
	classRepo := h.storageManager.ClassRepository()

	// Create the class/term directory
	ctx := r.Context()
	err := classRepo.CreateClassTerm(ctx, className)
	if err != nil {
		if os.IsExist(err) {
			HandleError(w, errors.Wrap(err, errors.ErrClassNotFound, "Class/term already exists"))
		} else if strings.Contains(err.Error(), "invalid class/term name") {
			HandleError(w, errors.Wrap(err, errors.ErrInvalidClassName, "Invalid class/term name"))
		} else {
			HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to create class/term directory"))
		}
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]string{"message": "Class/term '" + className + "' created successfully"})
}

// ListClassTermsHandler handles GET requests to /classes
// Response: Success (200) with list of class names or Error (500)
func (h *Handler) ListClassTermsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed"))
		return
	}

	// Get the class repository from the storage manager
	classRepo := h.storageManager.ClassRepository()

	// List all class/terms
	ctx := r.Context()
	classTerms, err := classRepo.ListClassTerms(ctx)
	if err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrInternal, "Failed to list class/terms"))
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
func (h *Handler) ProjectsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload projects
		var projects []string
		if err := json.NewDecoder(r.Body).Decode(&projects); err != nil {
			HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body, expected JSON array of project strings"))
			return
		}
		defer r.Body.Close()

		// Get the project repository from the storage manager
		projectRepo := h.storageManager.ProjectRepository()

		// Save the projects
		ctx := r.Context()
		err := projectRepo.SaveProjects(ctx, className, projects)
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				HandleError(w, errors.Wrap(err, errors.ErrClassNotFound, "Class/term not found"))
			} else {
				HandleError(w, errors.Wrap(err, errors.ErrFailedToSaveProjects, "Failed to save projects"))
			}
			return
		}
		RespondJSON(w, http.StatusCreated, map[string]string{"message": "Projects for '" + className + "' saved successfully"})

	case http.MethodGet: // Get projects
		// Get the project repository from the storage manager
		projectRepo := h.storageManager.ProjectRepository()

		// Load the projects
		ctx := r.Context()
		projects, err := projectRepo.LoadProjects(ctx, className)
		if err != nil {
			if os.IsNotExist(err) {
				HandleError(w, errors.Wrap(err, errors.ErrProjectsNotFound, "Projects not found"))
			} else {
				HandleError(w, errors.Wrap(err, errors.ErrFailedToLoadProjects, "Failed to load projects"))
			}
			return
		}
		if projects == nil { // Ensure empty list, not null
			projects = []string{}
		}
		RespondJSON(w, http.StatusOK, projects)

	default:
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET and POST methods are allowed for projects"))
	}
}

// StudentsHandler manages student-related requests for a specific class/term.
// It routes based on method (GET for loading, POST for saving).
// URL: /classes/{className}/students
func (h *Handler) StudentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	switch r.Method {
	case http.MethodPost: // Upload students
		var students []models.StudentInput // Type from models package
		if err := json.NewDecoder(r.Body).Decode(&students); err != nil {
			HandleError(w, errors.Wrap(err, errors.ErrInvalidRequestBody, "Invalid request body, expected JSON array of student objects"))
			return
		}
		defer r.Body.Close()

		// Get the student repository from the storage manager
		studentRepo := h.storageManager.StudentRepository()

		// Save the students
		ctx := r.Context()
		err := studentRepo.SaveStudents(ctx, className, students)
		if err != nil {
			if os.IsNotExist(err) { // Check if class directory doesn't exist
				HandleError(w, errors.Wrap(err, errors.ErrClassNotFound, "Class/term not found"))
			} else {
				HandleError(w, errors.Wrap(err, errors.ErrFailedToSaveStudents, "Failed to save students"))
			}
			return
		}
		RespondJSON(w, http.StatusCreated, map[string]string{"message": "Students for '" + className + "' saved successfully"})

	case http.MethodGet: // Get students
		// Get the student repository from the storage manager
		studentRepo := h.storageManager.StudentRepository()

		// Load the students
		ctx := r.Context()
		students, err := studentRepo.LoadStudents(ctx, className)
		if err != nil {
			if os.IsNotExist(err) {
				HandleError(w, errors.Wrap(err, errors.ErrStudentsNotFound, "Students not found"))
			} else {
				HandleError(w, errors.Wrap(err, errors.ErrFailedToLoadStudents, "Failed to load students"))
			}
			return
		}
		if students == nil { // Ensure empty list, not null
			students = []models.StudentInput{}
		}
		RespondJSON(w, http.StatusOK, students)

	default:
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET and POST methods are allowed for students"))
	}
}

// MasterClassResourceHandler is a generic handler for /classes/{className}/{resource}
// e.g., /classes/Fall2023/projects or /classes/Spring2024/students
// It will delegate to specific handlers like ProjectsHandler or StudentsHandler.
func (h *Handler) MasterClassResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/classes/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		// This case should ideally be handled by a more specific /classes handler if it's for listing or creating classes.
		// If it reaches here, it means it's an invalid path like /classes//projects
		HandleError(w, errors.New(errors.ErrInvalidRequestBody, "Class/term name missing in URL path"))
		return
	}
	className := parts[0]

	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	if len(parts) == 1 {
		// Potentially a handler for /classes/{className} - e.g., get class details (not specified yet)
		HandleError(w, errors.New(errors.ErrResourceNotFound, "Resource type missing in URL path"))
		return
	}

	resourceType := parts[1]
	switch resourceType {
	case "projects":
		h.ProjectsHandler(w, r, className)
	case "students":
		h.StudentsHandler(w, r, className)
	case "assign":
		// TriggerAssignmentHandler is defined in assignment_handler.go
		// It expects a POST request.
		h.TriggerAssignmentHandler(w, r, className) // This only handles POST
	case "assignments":
		// Path: /classes/{className}/assignments OR /classes/{className}/assignments/{assignmentID}
		if len(parts) == 2 { // Exactly /classes/{className}/assignments
			h.ListAssignmentsHandler(w, r, className) // Defined in assignment_handler.go
		} else if len(parts) == 3 { // Exactly /classes/{className}/assignments/{assignmentID}
			assignmentID := parts[2]
			if !h.IsValidClassName(assignmentID) { // Reuse validation for simplicity
				HandleError(w, errors.New(errors.ErrInvalidAssignmentID, "Invalid assignment ID format"))
				return
			}
			h.GetAssignmentResultHandler(w, r, className, assignmentID) // Defined in assignment_handler.go
		} else {
			// Path is too long, e.g., /classes/{className}/assignments/{id}/something_else
			HandleError(w, errors.New(errors.ErrResourceNotFound, "Invalid path structure for assignments resource"))
		}
	default:
		HandleError(w, errors.New(errors.ErrResourceNotFound, "Unknown resource type"))
	}
}

// HandleClassesBase routes requests for the /classes endpoint.
// POST to /classes -> CreateClassTermHandler
// GET to /classes  -> ListClassTermsHandler
func (h *Handler) HandleClassesBase(w http.ResponseWriter, r *http.Request) {
	// Ensure the path is exactly "/classes" and not "/classes/" or "/classes/something"
	if r.URL.Path != "/classes" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.CreateClassTermHandler(w, r)
	case http.MethodGet:
		h.ListClassTermsHandler(w, r)
	default:
		RespondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for /classes")
	}
}

// MasterRouter is the primary router for the application.
func (h *Handler) MasterRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/":
		// Handle root path if needed, e.g., a welcome message or API documentation link
		if r.Method == http.MethodGet {
			// Include version information in the welcome message
			RespondJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Welcome to the Fair Project Assignment API",
				"version": "1.0.0", // This should be updated to use the version package
			})
		} else {
			RespondError(w, http.StatusMethodNotAllowed, "Only GET is allowed for the root path")
		}
	case path == "/version":
		// Handle the /version endpoint
		VersionHandler(w, r)
	case path == "/docs":
		// Handle the /docs endpoint for API documentation
		http.NotFound(w, r) // This should be updated to use the docs package
	case path == "/auth/login" || strings.HasPrefix(path, "/auth/users") || path == "/user/whoami":
		// Handle authentication endpoints
		AuthResourceHandler(w, r)
	case path == "/classes":
		h.HandleClassesBase(w, r)
	case strings.HasPrefix(path, "/classes/"):
		// MasterClassResourceHandler handles paths like /classes/{className}/projects, /classes/{className}/students etc.
		h.MasterClassResourceHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// IsValidClassName checks if the className is valid for use as a directory name.
func (h *Handler) IsValidClassName(className string) bool {
	// Simple implementation using regex
	if className == "" {
		return false
	}

	// Use a simple regex check for now
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(className)
}
