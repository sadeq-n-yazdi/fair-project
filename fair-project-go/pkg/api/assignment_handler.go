package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sadeq/fair-project-go/pkg/assignment"
	"github.com/sadeq/fair-project-go/pkg/errors"
	"github.com/sadeq/fair-project-go/pkg/models"
)

// TriggerAssignmentHandler handles the process of assigning students to projects for a given class.
// POST /classes/{className}/assign
func (h *Handler) TriggerAssignmentHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only POST method is allowed for triggering assignments"))
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	// 1. Load projects
	projectRepo := h.storageManager.ProjectRepository()
	ctx := r.Context()
	projects, err := projectRepo.LoadProjects(ctx, className)
	if err != nil {
		if os.IsNotExist(err) {
			HandleError(w, errors.Wrap(err, errors.ErrProjectsNotFound, "Projects not found. Please upload projects first"))
		} else {
			HandleError(w, errors.Wrap(err, errors.ErrFailedToLoadProjects, "Failed to load projects"))
		}
		return
	}
	if len(projects) == 0 {
		HandleError(w, errors.New(errors.ErrNoProjects, "No projects found. Assignment cannot run without projects"))
		return
	}

	// 2. Load students input data
	studentRepo := h.storageManager.StudentRepository()
	studentInputs, err := studentRepo.LoadStudents(ctx, className)
	if err != nil {
		if os.IsNotExist(err) {
			HandleError(w, errors.Wrap(err, errors.ErrStudentsNotFound, "Students data not found. Please upload students data first"))
		} else {
			HandleError(w, errors.Wrap(err, errors.ErrFailedToLoadStudents, "Failed to load students data"))
		}
		return
	}
	if len(studentInputs) == 0 {
		HandleError(w, errors.New(errors.ErrNoStudents, "No students found. Assignment cannot run without students"))
		return
	}

	// 3. Transform []StudentInput to []Student
	var students []models.Student // This is the type expected by assignProjects
	for _, si := range studentInputs {
		students = append(students, models.Student{
			Name:            si.Name,
			Preferences:     si.Preferences, // Assuming preferences are directly usable
			Assigned:        false,
			AssignedProject: "",
			BadPreferences:  []string{}, // Initialize as empty slice
		})
	}

	// 4. Call the core assignment logic
	selectedByProjects, selectedByStudent, updatedStudents, remainingProjects := assignment.AssignProjects(students, projects)

	// 5. Create AssignmentOutput
	assignmentOutputData := models.AssignmentOutput{
		SelectedByProjects: selectedByProjects,
		SelectedByStudent:  selectedByStudent,
		Students:           updatedStudents,
		RemainingProjects:  remainingProjects,
	}

	// 6. Generate a unique assignmentID (timestamp-based for now)
	assignmentID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 7. Save the results
	assignmentRepo := h.storageManager.AssignmentRepository()
	if err := assignmentRepo.SaveAssignmentResults(ctx, className, assignmentID, assignmentOutputData); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrFailedToSaveAssignment, "Failed to save assignment results"))
		return
	}

	// 8. Respond with success and the assignment output
	responsePayload := map[string]interface{}{
		"message":      "Assignment completed successfully",
		"class":        className,
		"assignmentId": assignmentID,
		"results":      assignmentOutputData,
	}
	RespondJSON(w, http.StatusCreated, responsePayload)
}

// ListAssignmentsHandler handles GET requests to list all assignment IDs for a class.
// GET /classes/{className}/assignments
func (h *Handler) ListAssignmentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed for listing assignments"))
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	// Get the assignment repository from the storage manager
	assignmentRepo := h.storageManager.AssignmentRepository()

	// List all assignment results
	ctx := r.Context()
	assignmentIDs, err := assignmentRepo.ListAssignmentResults(ctx, className)
	if err != nil {
		if os.IsNotExist(err) { // Directory for className might not exist
			HandleError(w, errors.Wrap(err, errors.ErrClassNotFound, "Class/term not found or has no assignments"))
		} else {
			HandleError(w, errors.Wrap(err, errors.ErrFailedToListAssignments, "Failed to list assignments"))
		}
		return
	}

	if assignmentIDs == nil { // Ensure empty list, not null
		assignmentIDs = []string{}
	}
	RespondJSON(w, http.StatusOK, assignmentIDs)
}

// GetAssignmentResultHandler handles GET requests for a specific assignment result.
// GET /classes/{className}/assignments/{assignmentID}
func (h *Handler) GetAssignmentResultHandler(w http.ResponseWriter, r *http.Request, className string, assignmentID string) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed for fetching an assignment result"))
		return
	}

	// Validate className and assignmentID (already done by masterClassResourceHandler, but good for defense)
	if !h.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}
	if !h.IsValidClassName(assignmentID) { // Using same validation as className for simplicity
		HandleError(w, errors.New(errors.ErrInvalidAssignmentID, "Invalid assignment ID format"))
		return
	}

	// Get the assignment repository from the storage manager
	assignmentRepo := h.storageManager.AssignmentRepository()

	// Load the assignment result
	ctx := r.Context()
	result, err := assignmentRepo.LoadAssignmentResult(ctx, className, assignmentID)
	if err != nil {
		if os.IsNotExist(err) {
			HandleError(w, errors.Wrap(err, errors.ErrAssignmentNotFound, "Assignment result not found"))
		} else {
			HandleError(w, errors.Wrap(err, errors.ErrFailedToLoadAssignment, "Failed to load assignment result"))
		}
		return
	}
	RespondJSON(w, http.StatusOK, result)
}
