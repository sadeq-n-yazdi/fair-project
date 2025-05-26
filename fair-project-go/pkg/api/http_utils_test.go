package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Create a payload
	payload := map[string]string{"message": "test message"}
	
	// Call the function
	RespondJSON(rr, http.StatusOK, payload)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
	
	// Check the response body
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}
	
	if response["message"] != "test message" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["message"], "test message")
	}
}

func TestRespondError(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the function
	RespondError(rr, http.StatusBadRequest, "test error")
	
	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
	
	// Check the response body
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}
	
	if response["error"] != "test error" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["error"], "test error")
	}
}
