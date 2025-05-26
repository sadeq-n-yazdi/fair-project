package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/sadeq/fair-project-go/pkg/api"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// handleClassesBase routes requests for the /classes endpoint.
// POST to /classes -> CreateClassTermHandler
// GET to /classes  -> ListClassTermsHandler
func handleClassesBase(w http.ResponseWriter, r *http.Request) {
	// Ensure the path is exactly "/classes" and not "/classes/" or "/classes/something"
	// The check for "/classes/" is implicitly handled by registering MasterClassResourceHandler
	// for "/classes/" which is a more specific match for paths with a trailing slash.
	// However, adding an explicit check here ensures this handler only manages "/classes".
	if r.URL.Path != "/classes" {
		// This situation should ideally not occur if routing is set up correctly with
		// MasterClassResourceHandler for "/classes/".
		// If it does, it implies a request like "/classesX" or some other unexpected path.
		// However, the default mux behavior might route "/classesanything" here if not for
		// the more specific "/classes/" route. For strictness:
		http.NotFound(w, r) // Or RespondError with a specific message
		return
	}

	switch r.Method {
	case http.MethodPost:
		api.CreateClassTermHandler(w, r) // Defined in class_handlers.go
	case http.MethodGet:
		api.ListClassTermsHandler(w, r) // Defined in class_handlers.go
	default:
		// RespondError is defined in http_utils.go
		api.RespondError(w, http.StatusMethodNotAllowed, "Only GET and POST methods are allowed for /classes")
	}
}

// masterRouter is the primary router for the application.
// It replaces the previous MasterClassResourceHandler to also handle the root path
// and provide clearer separation of concerns.
func masterRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/":
		// Handle root path if needed, e.g., a welcome message or API documentation link
		if r.Method == http.MethodGet {
			api.RespondJSON(w, http.StatusOK, map[string]string{"message": "Welcome to the Fair Project Assignment API"})
		} else {
			api.RespondError(w, http.StatusMethodNotAllowed, "Only GET is allowed for the root path")
		}
	case path == "/classes":
		handleClassesBase(w, r)
	case strings.HasPrefix(path, "/classes/"):
		// MasterClassResourceHandler is defined in class_handlers.go
		// It handles paths like /classes/{className}/projects, /classes/{className}/students etc.
		api.MasterClassResourceHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	// 1. Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// 2. Ensure the base data directory exists
	if err := storage.EnsureBaseDir(); err != nil { // From storage.go
		log.Fatalf("Failed to ensure base data directory: %v", err)
	}

	// 3. Register handlers using the masterRouter
	// The masterRouter will delegate to the appropriate handlers.
	http.HandleFunc("/", masterRouter)

	// 4. Start the HTTP server
	port := ":8080"
	log.Printf("Starting server on port %s. Listening for requests on / , /classes, and /classes/...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
