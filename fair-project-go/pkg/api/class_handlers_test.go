package api

import (
	"net/http"
	"testing"

	"github.com/sadeq/fair-project-go/pkg/models"
)

func TestCreateClassTermHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid class name",
			requestBody:    map[string]string{"name": "TestClass"},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "Empty class name",
			requestBody:    map[string]string{"name": ""},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid class name",
			requestBody:    map[string]string{"name": "Test Class"}, // Contains space
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := createTestRequest(t, http.MethodPost, "/classes", tc.requestBody)

			// Execute request
			rr := executeRequest(req, CreateClassTermHandler)

			// Check response
			checkResponseCode(t, tc.expectedStatus, rr.Code)
			checkContentType(t, rr, "application/json")

			// Parse response
			if tc.expectedError {
				// For error responses, expect a nested error object
				var response map[string]map[string]string
				parseResponse(t, rr, &response)

				// Check response content
				if errorObj, ok := response["error"]; !ok {
					t.Errorf("Expected error object in response, got: %v", response)
				} else if _, ok := errorObj["code"]; !ok {
					t.Errorf("Expected error code in response, got: %v", response)
				} else if _, ok := errorObj["message"]; !ok {
					t.Errorf("Expected error message in response, got: %v", response)
				}
			} else {
				// For success responses, expect a simple message
				var response map[string]string
				parseResponse(t, rr, &response)

				if _, ok := response["message"]; !ok {
					t.Errorf("Expected message in response, got: %v", response)
				}
			}
		})
	}
}

func TestListClassTermsHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create some test classes
	testClasses := []string{"Class1", "Class2", "Class3"}
	for _, className := range testClasses {
		createTestClass(t, className)
	}

	// Create request
	req := createTestRequest(t, http.MethodGet, "/classes", nil)

	// Execute request
	rr := executeRequest(req, ListClassTermsHandler)

	// Check response
	checkResponseCode(t, http.StatusOK, rr.Code)
	checkContentType(t, rr, "application/json")

	// Parse response
	var response []string
	parseResponse(t, rr, &response)

	// Check response content
	if len(response) != len(testClasses) {
		t.Errorf("Expected %d classes, got %d", len(testClasses), len(response))
	}

	// Check that all test classes are in the response
	for _, className := range testClasses {
		found := false
		for _, respClass := range response {
			if respClass == className {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected class %s in response, but not found", className)
		}
	}
}

func TestProjectsHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test class
	className := "TestClass"
	createTestClass(t, className)

	// Test POST /classes/{className}/projects
	t.Run("POST projects", func(t *testing.T) {
		// Test data
		projects := []string{"Project1", "Project2", "Project3"}

		// Create request
		req := createTestRequest(t, http.MethodPost, "/classes/"+className+"/projects", projects)

		// Create a handler function that calls ProjectsHandler with the className
		handler := func(w http.ResponseWriter, r *http.Request) {
			ProjectsHandler(w, r, className)
		}

		// Execute request
		rr := executeRequest(req, handler)

		// Check response
		checkResponseCode(t, http.StatusCreated, rr.Code)
		checkContentType(t, rr, "application/json")

		// Parse response
		var response map[string]string
		parseResponse(t, rr, &response)

		// Check response content
		if _, ok := response["message"]; !ok {
			t.Errorf("Expected message in response, got: %v", response)
		}
	})

	// Test GET /classes/{className}/projects
	t.Run("GET projects", func(t *testing.T) {
		// Create some test projects
		projects := []string{"Project1", "Project2", "Project3"}
		createTestProjects(t, className, projects)

		// Create request
		req := createTestRequest(t, http.MethodGet, "/classes/"+className+"/projects", nil)

		// Create a handler function that calls ProjectsHandler with the className
		handler := func(w http.ResponseWriter, r *http.Request) {
			ProjectsHandler(w, r, className)
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
		if len(response) != len(projects) {
			t.Errorf("Expected %d projects, got %d", len(projects), len(response))
		}

		// Check that all test projects are in the response
		for _, project := range projects {
			found := false
			for _, respProject := range response {
				if respProject == project {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected project %s in response, but not found", project)
			}
		}
	})
}

func TestStudentsHandler(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test class
	className := "TestClass"
	createTestClass(t, className)

	// Test POST /classes/{className}/students
	t.Run("POST students", func(t *testing.T) {
		// Test data
		students := []models.StudentInput{
			{Name: "Student1", Preferences: []string{"Project1", "Project2"}},
			{Name: "Student2", Preferences: []string{"Project2", "Project3"}},
		}

		// Create request
		req := createTestRequest(t, http.MethodPost, "/classes/"+className+"/students", students)

		// Create a handler function that calls StudentsHandler with the className
		handler := func(w http.ResponseWriter, r *http.Request) {
			StudentsHandler(w, r, className)
		}

		// Execute request
		rr := executeRequest(req, handler)

		// Check response
		checkResponseCode(t, http.StatusCreated, rr.Code)
		checkContentType(t, rr, "application/json")

		// Parse response
		var response map[string]string
		parseResponse(t, rr, &response)

		// Check response content
		if _, ok := response["message"]; !ok {
			t.Errorf("Expected message in response, got: %v", response)
		}
	})

	// Test GET /classes/{className}/students
	t.Run("GET students", func(t *testing.T) {
		// Create some test students
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
		req := createTestRequest(t, http.MethodGet, "/classes/"+className+"/students", nil)

		// Create a handler function that calls StudentsHandler with the className
		handler := func(w http.ResponseWriter, r *http.Request) {
			StudentsHandler(w, r, className)
		}

		// Execute request
		rr := executeRequest(req, handler)

		// Check response
		checkResponseCode(t, http.StatusOK, rr.Code)
		checkContentType(t, rr, "application/json")

		// Parse response
		var response []models.StudentInput
		parseResponse(t, rr, &response)

		// Check response content
		if len(response) != len(students) {
			t.Errorf("Expected %d students, got %d", len(students), len(response))
		}
	})
}
