package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	// Ensure other necessary imports from the project are implicitly available
	// e.g., types from types.go, storage functions from storage.go, etc.
)

// triggerAssignmentHandler handles the process of assigning students to projects for a given class.
// POST /classes/{className}/assign
func triggerAssignmentHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Only POST method is allowed for triggering assignments")
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL: "+className)
		return
	}

	// 1. Load projects
	projects, err := LoadProjects(className) // From storage.go
	if err != nil {
		if os.IsNotExist(err) {
			respondError(w, http.StatusNotFound, "Projects not found for class/term '"+className+"'. Please upload projects first.")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to load projects: "+err.Error())
		}
		return
	}
	if len(projects) == 0 {
		respondError(w, http.StatusUnprocessableEntity, "No projects found for class/term '"+className+"'. Assignment cannot run without projects.")
		return
	}

	// 2. Load students input data
	studentInputs, err := LoadStudents(className) // From storage.go
	if err != nil {
		if os.IsNotExist(err) {
			respondError(w, http.StatusNotFound, "Students data not found for class/term '"+className+"'. Please upload students data first.")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to load students data: "+err.Error())
		}
		return
	}
	if len(studentInputs) == 0 {
		respondError(w, http.StatusUnprocessableEntity, "No students found for class/term '"+className+"'. Assignment cannot run without students.")
		return
	}

	// 3. Transform []StudentInput to []Student
	var students []Student // This is the type expected by assignProjects
	for _, si := range studentInputs {
		students = append(students, Student{
			Name:            si.Name,
			Preferences:     si.Preferences, // Assuming preferences are directly usable
			Assigned:        false,
			AssignedProject: "",
			BadPreferences:  []string{}, // Initialize as empty slice
		})
	}

	// 4. Call the core assignment logic
	// assignProjects is from assignment.go
	selectedByProjects, selectedByStudent, updatedStudents, remainingProjects := assignProjects(students, projects)

	// 5. Create AssignmentOutput
	assignmentOutputData := AssignmentOutput{
		SelectedByProjects: selectedByProjects,
		SelectedByStudent:  selectedByStudent,
		Students:           updatedStudents,
		RemainingProjects:  remainingProjects,
	}

	// 6. Generate a unique assignmentID (timestamp-based for now)
	assignmentID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 7. Save the results
	// SaveAssignmentResults is from storage.go
	if err := SaveAssignmentResults(className, assignmentID, assignmentOutputData); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save assignment results: "+err.Error())
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
	respondJSON(w, http.StatusOK, responsePayload)
}

// listAssignmentsHandler handles GET requests to list all assignment IDs for a class.
// GET /classes/{className}/assignments
func listAssignmentsHandler(w http.ResponseWriter, r *http.Request, className string) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed for listing assignments")
		return
	}

	// Validate className (already done by masterClassResourceHandler, but good for defense)
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL: "+className)
		return
	}

	assignmentIDs, err := ListAssignmentResults(className) // From storage.go
	if err != nil {
		if os.IsNotExist(err) { // Directory for className might not exist
			respondError(w, http.StatusNotFound, "Class/term '"+className+"' not found or has no assignments.")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to list assignments: "+err.Error())
		}
		return
	}

	if assignmentIDs == nil { // Ensure empty list, not null
		assignmentIDs = []string{}
	}
	respondJSON(w, http.StatusOK, assignmentIDs)
}

// getAssignmentResultHandler handles GET requests for a specific assignment result.
// GET /classes/{className}/assignments/{assignmentID}
func getAssignmentResultHandler(w http.ResponseWriter, r *http.Request, className string, assignmentID string) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed for fetching an assignment result")
		return
	}

	// Validate className and assignmentID (already done by masterClassResourceHandler, but good for defense)
	if !isValidClassName(className) {
		respondError(w, http.StatusBadRequest, "Invalid class/term name format in URL: "+className)
		return
	}
	if !isValidClassName(assignmentID) { // Using same validation as className for simplicity
		respondError(w, http.StatusBadRequest, "Invalid assignment ID format in URL: "+assignmentID)
		return
	}

	result, err := LoadAssignmentResult(className, assignmentID) // From storage.go
	if err != nil {
		if os.IsNotExist(err) {
			respondError(w, http.StatusNotFound, "Assignment result '"+assignmentID+"' not found for class/term '"+className+"'.")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to load assignment result: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, result)
}
