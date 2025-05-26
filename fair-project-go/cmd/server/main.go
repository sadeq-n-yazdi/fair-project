package main

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sadeq/fair-project-go/pkg/api"
	authpkg "github.com/sadeq/fair-project-go/pkg/auth"
	"github.com/sadeq/fair-project-go/pkg/config"
	"github.com/sadeq/fair-project-go/pkg/docs"
	authmiddleware "github.com/sadeq/fair-project-go/pkg/middleware/auth"
	"github.com/sadeq/fair-project-go/pkg/middleware/logging"
	"github.com/sadeq/fair-project-go/pkg/storage"
	"github.com/sadeq/fair-project-go/pkg/version"
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
			// Include version information in the welcome message
			api.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Welcome to the Fair Project Assignment API",
				"version": version.String(),
			})
		} else {
			api.RespondError(w, http.StatusMethodNotAllowed, "Only GET is allowed for the root path")
		}
	case path == "/version":
		// Handle the /version endpoint
		api.VersionHandler(w, r)
	case path == "/docs":
		// Handle the /docs endpoint for API documentation
		docs.Handler(w, r)
	case path == "/auth/login" || strings.HasPrefix(path, "/auth/users"):
		// Handle authentication endpoints
		api.AuthResourceHandler(w, r)
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

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	// Try to open the .env file
	file, err := os.Open(".env")
	if err != nil {
		// Try to find the .env file in the parent directory
		file, err = os.Open("../.env")
		if err != nil {
			log.Printf("Warning: Could not open .env file: %v", err)
			return
		}
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Warning: Error reading .env file: %v", err)
	}
}

func main() {
	// 1. Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// 1.5. Load environment variables from .env file
	loadEnvFile()

	// 2. Initialize configuration from environment variables
	config.InitFromEnv()
	log.Printf("Log level set to: %s", config.GetLogLevel())

	// 3. Set the base data directory from environment variable or use default
	dataDir := os.Getenv("DATA_DIR")
	if dataDir != "" {
		storage.SetBaseDataDir(dataDir)
		log.Printf("Using data directory: %s", dataDir)
	} else {
		log.Printf("Using default data directory: %s", storage.GetBaseDataDir())
	}

	// 4. Ensure the base data directory exists
	if err := storage.EnsureBaseDir(); err != nil { // From storage.go
		log.Fatalf("Failed to ensure base data directory: %v", err)
	}

	// 5. Initialize default users
	if err := storage.InitializeDefaultUsers(); err != nil {
		log.Fatalf("Failed to initialize default users: %v", err)
	}

	// 6. Set JWT secret from environment variable or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		authpkg.SetJWTSecret(jwtSecret)
	}

	// 7. Create a handler with the logging and authentication middleware
	handler := logging.Middleware(authmiddleware.Middleware(http.HandlerFunc(masterRouter)))

	// 8. Get the port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	// 9. Start the HTTP server with the middleware
	log.Printf("Starting server on port %s. Listening for requests on / , /classes, and /classes/...\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
