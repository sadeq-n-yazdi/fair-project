package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sadeq/fair-project-go/pkg/assignment"
	"github.com/sadeq/fair-project-go/pkg/errors"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// TriggerAssignmentHandler handles the process of assigning students to projects for a given class.
// POST /classes/{className}/assign
func TriggerAssignmentHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodPost {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only POST method is allowed for triggering assignments"))
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !storage.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	// 1. Load projects
	projects, err := storage.LoadProjects(className) // From storage.go
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
	studentInputs, err := storage.LoadStudents(className) // From storage.go
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
	// assignProjects is from assignment.go
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
	// SaveAssignmentResults is from storage.go
	if err := storage.SaveAssignmentResults(className, assignmentID, assignmentOutputData); err != nil {
		HandleError(w, errors.Wrap(err, errors.ErrFailedToSaveAssignment, "Failed to save assignment results"))
		return
	}

	// 8. Respond with success and the assignment output
	// Adding assignmentId to the response along with the full output might be redundant if the output itself is keyed by ID in some contexts,
	// but providing it directly can be convenient for clients.
	// For now, just return the AssignmentOutput as requested.
	// To include assignmentID in the response body along with the results, we'd need to wrap them.
	// Let's stick to returning the AssignmentOutput directly.
	// A 201 Created might be more appropriate if the assignment itself is seen as a new resource.
	// However, since we save it and give it an ID, and the request is to "trigger" an action,
	// 200 OK with the result is also common. Let's use 200 OK.

	// To include assignmentID and class in the response, we can create a wrapper struct or a map
	responsePayload := map[string]interface{}{
		"message":      "Assignment completed successfully",
		"class":        className,
		"assignmentId": assignmentID,
		"results":      assignmentOutputData,
	}
	RespondJSON(w, http.StatusOK, responsePayload)
}

// ListAssignmentsHandler handles GET requests to list all assignment IDs for a class.
// GET /classes/{className}/assignments
func ListAssignmentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed for listing assignments"))
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !storage.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}

	assignmentIDs, err := storage.ListAssignmentResults(className) // From storage.go
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
func GetAssignmentResultHandler(w http.ResponseWriter, r *http.Request, className string, assignmentID string) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed for fetching an assignment result"))
		return
	}

	// Validate className and assignmentID (already done by masterClassResourceHandler, but good for defense)
	if !storage.IsValidClassName(className) {
		HandleError(w, errors.New(errors.ErrInvalidClassName, "Invalid class/term name format"))
		return
	}
	if !storage.IsValidClassName(assignmentID) { // Using same validation as className for simplicity
		HandleError(w, errors.New(errors.ErrInvalidAssignmentID, "Invalid assignment ID format"))
		return
	}

	result, err := storage.LoadAssignmentResult(className, assignmentID) // From storage.go
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
