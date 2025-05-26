package api

import (
	"net/http"
	"testing"

	"github.com/sadeq/fair-project-go/pkg/models"
)

func TestTriggerAssignmentHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test class
	className := "TestClass"
	createTestClass(t, className)

	// Create test projects
	projects := []string{"Project1", "Project2", "Project3"}
	createTestProjects(t, className, projects)

	// Create test students
	students := []interface{}{
		map[string]interface{}{
			"name":        "Student1",
			"preferences": []string{"Project1", "Project2"},
		},
		map[string]interface{}{
			"name":        "Student2",
			"preferences": []string{"Project2", "Project3"},
		},
	}
	createTestStudents(t, className, students)

	// Create request
	req := createTestRequest(t, http.MethodPost, "/classes/"+className+"/assign", nil)

	// Create a handler function that calls TriggerAssignmentHandler with the className
	handler := func(w http.ResponseWriter, r *http.Request) {
		TriggerAssignmentHandler(w, r, className)
	}

	// Execute request
	rr := executeRequest(req, handler)

	// Check response
	checkResponseCode(t, http.StatusCreated, rr.Code)
	checkContentType(t, rr, "application/json")

	// Parse response
	var response map[string]interface{}
	parseResponse(t, rr, &response)

	// Check response content
	if _, ok := response["message"]; !ok {
		t.Errorf("Expected message in response, got: %v", response)
	}
	if _, ok := response["assignmentId"]; !ok {
		t.Errorf("Expected assignmentId in response, got: %v", response)
	}
	if _, ok := response["results"]; !ok {
		t.Errorf("Expected results in response, got: %v", response)
	}
}

func TestListAssignmentsHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test class
	className := "TestClass"
	createTestClass(t, className)

	// Create test projects
	projects := []string{"Project1", "Project2", "Project3"}
	createTestProjects(t, className, projects)

	// Create test students
	students := []interface{}{
		map[string]interface{}{
			"name":        "Student1",
			"preferences": []string{"Project1", "Project2"},
		},
		map[string]interface{}{
			"name":        "Student2",
			"preferences": []string{"Project2", "Project3"},
		},
	}
	createTestStudents(t, className, students)

	// Trigger an assignment to create an assignment result
	triggerReq := createTestRequest(t, http.MethodPost, "/classes/"+className+"/assign", nil)
	triggerHandler := func(w http.ResponseWriter, r *http.Request) {
		TriggerAssignmentHandler(w, r, className)
	}
	executeRequest(triggerReq, triggerHandler)

	// Create request to list assignments
	req := createTestRequest(t, http.MethodGet, "/classes/"+className+"/assignments", nil)

	// Create a handler function that calls ListAssignmentsHandler with the className
	handler := func(w http.ResponseWriter, r *http.Request) {
		ListAssignmentsHandler(w, r, className)
	}

	// Execute request
	rr := executeRequest(req, handler)

	// Check response
	checkResponseCode(t, http.StatusOK, rr.Code)
	checkContentType(t, rr, "application/json")

	// Parse response
	var response []string
	parseResponse(t, rr, &response)

	// Check response content
	if len(response) != 1 {
		t.Errorf("Expected 1 assignment, got %d", len(response))
	}
}

func TestGetAssignmentResultHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test class
	className := "TestClass"
	createTestClass(t, className)

	// Create test projects
	projects := []string{"Project1", "Project2", "Project3"}
	createTestProjects(t, className, projects)

	// Create test students
	students := []interface{}{
		map[string]interface{}{
			"name":        "Student1",
			"preferences": []string{"Project1", "Project2"},
		},
		map[string]interface{}{
			"name":        "Student2",
			"preferences": []string{"Project2", "Project3"},
		},
	}
	createTestStudents(t, className, students)

	// Trigger an assignment to create an assignment result
	triggerReq := createTestRequest(t, http.MethodPost, "/classes/"+className+"/assign", nil)
	triggerHandler := func(w http.ResponseWriter, r *http.Request) {
		TriggerAssignmentHandler(w, r, className)
	}
	triggerRR := executeRequest(triggerReq, triggerHandler)

	// Get the assignment ID from the response
	var triggerResponse map[string]interface{}
	parseResponse(t, triggerRR, &triggerResponse)
	assignmentID := triggerResponse["assignmentId"].(string)

	// Create request to get assignment result
	req := createTestRequest(t, http.MethodGet, "/classes/"+className+"/assignments/"+assignmentID, nil)

	// Create a handler function that calls GetAssignmentResultHandler with the className and assignmentID
	handler := func(w http.ResponseWriter, r *http.Request) {
		GetAssignmentResultHandler(w, r, className, assignmentID)
	}

	// Execute request
	rr := executeRequest(req, handler)

	// Check response
	checkResponseCode(t, http.StatusOK, rr.Code)
	checkContentType(t, rr, "application/json")

	// Parse response
	var response models.AssignmentOutput
	parseResponse(t, rr, &response)

	// Check response content
	if len(response.Students) != 2 {
		t.Errorf("Expected 2 students in assignment result, got %d", len(response.Students))
	}
}
