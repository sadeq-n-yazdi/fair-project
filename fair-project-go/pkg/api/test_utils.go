package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sadeq/fair-project-go/pkg/storage"
)

// setupTestEnvironment creates a temporary directory for testing
// and sets up the storage package to use it.
// It returns a cleanup function that should be deferred.
func setupTestEnvironment(t *testing.T) func() {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fair-project-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override the baseDataDir in the storage package
	originalBaseDataDir := storage.GetBaseDataDir()
	storage.SetBaseDataDir(tempDir)

	// Return a cleanup function
	return func() {
		// Restore the original baseDataDir
		storage.SetBaseDataDir(originalBaseDataDir)

		// Remove the temporary directory
		os.RemoveAll(tempDir)
	}
}

// createTestRequest creates an HTTP request with the given method, URL, and body.
func createTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}

// executeRequest executes the given request and returns a response recorder.
func executeRequest(req *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

// checkResponseCode checks if the response code is as expected.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

// checkContentType checks if the content type is as expected.
func checkContentType(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("Expected content type %s. Got %s", expected, contentType)
	}
}

// parseResponse parses the response body into the given target.
func parseResponse(t *testing.T, rr *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(rr.Body.Bytes(), target)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}

// createTestClass creates a test class for testing.
func createTestClass(t *testing.T, className string) {
	err := storage.CreateClassTermDir(className)
	if err != nil {
		t.Fatalf("Failed to create test class: %v", err)
	}
}

// createTestProjects creates test projects for a class.
func createTestProjects(t *testing.T, className string, projects []string) {
	err := storage.SaveProjects(className, projects)
	if err != nil {
		t.Fatalf("Failed to create test projects: %v", err)
	}
}

// createTestStudents creates test students for a class.
func createTestStudents(t *testing.T, className string, students []interface{}) {
	// Convert []interface{} to []models.StudentInput
	studentsJSON, err := json.Marshal(students)
	if err != nil {
		t.Fatalf("Failed to marshal students: %v", err)
	}

	// Create the class directory if it doesn't exist
	classDir := filepath.Join(storage.GetBaseDataDir(), className)
	if _, err := os.Stat(classDir); os.IsNotExist(err) {
		err = os.MkdirAll(classDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create class directory: %v", err)
		}
	}

	// Write the students JSON file
	studentsFile := filepath.Join(classDir, "students.json")
	err = os.WriteFile(studentsFile, studentsJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write students file: %v", err)
	}
}
