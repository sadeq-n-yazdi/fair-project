package main

import (
	"encoding/json"
	"net/http"
)

// respondJSON sends a JSON response with the given status code and payload.
// It marshals the payload to JSON with indentation for readability.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ") // Indent for readability
	if err != nil {
		// If marshalling fails, log the error and send a generic server error
		// In a real app, you might want more sophisticated error logging here
		http.Error(w, "Error preparing response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response) // response is []byte
}

// respondError sends a JSON error message with the given status code.
// The error message is wrapped in a map: {"error": "message"}.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
