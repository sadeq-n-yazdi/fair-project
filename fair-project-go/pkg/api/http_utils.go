package api

import (
	"encoding/json"
	"net/http"

	"github.com/sadeq/fair-project-go/pkg/errors"
)

// RespondJSON sends a JSON response with the given status code and payload.
// It marshals the payload to JSON with indentation for readability.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ") // Indent for readability
	if err != nil {
		// If marshalling fails, log the error and send a generic server error
		appErr := errors.Wrap(err, errors.ErrInternal, "Error preparing response")
		errors.ToHTTPResponse(w, appErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response) // response is []byte
}

// RespondError sends a JSON error message with the given status code.
// This is a legacy function that should be replaced with errors.ToHTTPResponse.
// It's kept for backward compatibility.
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

// HandleError converts an error to an HTTP response using the errors package.
// This is the preferred way to handle errors in the API.
func HandleError(w http.ResponseWriter, err error) {
	errors.ToHTTPResponse(w, err)
}
