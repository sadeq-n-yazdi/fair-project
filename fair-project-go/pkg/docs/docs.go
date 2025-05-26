package docs

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

//go:embed openapi.json
var openAPISpec []byte

//go:embed swagger-ui.html
var swaggerUITemplate string

// Handler handles requests to the /docs endpoint.
// It serves the OpenAPI specification as JSON when requested with Accept: application/json
// or when the format query parameter is set to "json".
// Otherwise, it serves a web UI (Swagger UI) for the documentation.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if JSON format is requested
	format := r.URL.Query().Get("format")
	acceptHeader := r.Header.Get("Accept")

	if format == "json" || strings.Contains(acceptHeader, "application/json") {
		// Serve the OpenAPI specification as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write(openAPISpec)
		return
	}

	// Parse the OpenAPI spec to get the title
	var spec map[string]interface{}
	if err := json.Unmarshal(openAPISpec, &spec); err != nil {
		http.Error(w, "Error parsing OpenAPI specification", http.StatusInternalServerError)
		return
	}

	// Extract API title from the spec
	info, ok := spec["info"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid OpenAPI specification", http.StatusInternalServerError)
		return
	}
	title, ok := info["title"].(string)
	if !ok {
		title = "API Documentation"
	}

	// Serve the Swagger UI
	tmpl, err := template.New("swagger-ui").Parse(swaggerUITemplate)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	data := struct {
		Title    string
		SpecJSON string
	}{
		Title:    title,
		SpecJSON: string(openAPISpec),
	}
	tmpl.Execute(w, data)
}
