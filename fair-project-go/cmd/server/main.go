package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	case path == "/auth/login" || strings.HasPrefix(path, "/auth/users") || path == "/user/whoami":
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

// handleCompletion generates shell completion scripts
func handleCompletion(shell, outFile string) error {
	if shell == "" {
		return fmt.Errorf("shell type is required")
	}

	// Normalize shell name
	shell = strings.ToLower(shell)

	// Get the binary name
	binaryName := filepath.Base(os.Args[0])

	// Create the completion script
	var script string
	switch shell {
	case "bash":
		script = generateBashCompletion(binaryName)
	case "zsh":
		script = generateZshCompletion(binaryName)
	case "fish":
		script = generateFishCompletion(binaryName)
	default:
		return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh, fish)", shell)
	}

	// Write to file or stdout
	if outFile != "" {
		if err := os.WriteFile(outFile, []byte(script), 0644); err != nil {
			return fmt.Errorf("failed to write completion script to %s: %w", outFile, err)
		}
		fmt.Printf("Completion script written to %s\n", outFile)
		fmt.Printf("To install, run:\n")
		switch shell {
		case "bash":
			fmt.Printf("  echo \"source %s\" >> ~/.bashrc\n", outFile)
		case "zsh":
			fmt.Printf("  echo \"source %s\" >> ~/.zshrc\n", outFile)
		case "fish":
			fmt.Printf("  echo \"source %s\" >> ~/.config/fish/config.fish\n", outFile)
		}
	} else {
		fmt.Println(script)
		fmt.Printf("\nTo install, add the above to your shell configuration file or run:\n")
		switch shell {
		case "bash":
			fmt.Printf("  %s --completion bash > ~/.%s-completion.bash && echo \"source ~/.%s-completion.bash\" >> ~/.bashrc\n",
				binaryName, binaryName, binaryName)
		case "zsh":
			fmt.Printf("  %s --completion zsh > ~/.%s-completion.zsh && echo \"source ~/.%s-completion.zsh\" >> ~/.zshrc\n",
				binaryName, binaryName, binaryName)
		case "fish":
			fmt.Printf("  %s --completion fish > ~/.config/fish/%s-completion.fish\n", binaryName, binaryName)
		}
	}

	return nil
}

// generateBashCompletion generates a bash completion script
func generateBashCompletion(binaryName string) string {
	return fmt.Sprintf(`#!/bin/bash

_%s_completions() {
  COMPREPLY=()
  local word="${COMP_WORDS[COMP_CWORD]}"
  local completions="$(COMP_LINE="${COMP_LINE}" COMP_POINT="${COMP_POINT}" %s __complete)"
  COMPREPLY=( $(compgen -W "$completions" -- "$word") )
}

complete -F _%s_completions %s
`, binaryName, binaryName, binaryName, binaryName)
}

// generateZshCompletion generates a zsh completion script
func generateZshCompletion(binaryName string) string {
	return fmt.Sprintf(`#compdef %s

_%s() {
  local -a completions
  completions=("${(@f)$(COMP_LINE="${words[*]}" COMP_POINT=$#words %s __complete)}")
  _describe 'completions' completions
}

compdef _%s %s
`, binaryName, binaryName, binaryName, binaryName, binaryName)
}

// generateFishCompletion generates a fish completion script
func generateFishCompletion(binaryName string) string {
	return fmt.Sprintf(`function __fish_%s_complete
  set -l cl (commandline --tokenize --current-process)
  set -l tokens (commandline --tokenize --cut-at-cursor --current-process)
  %s __complete $tokens | tr '\n' ' '
end

complete -f -c %s -a '(__fish_%s_complete)'
`, binaryName, binaryName, binaryName, binaryName)
}

// handleCompletionRequest handles the __complete command for shell completion
func handleCompletionRequest() error {
	// Get the completion line from the environment
	line := os.Getenv("COMP_LINE")
	if line == "" {
		return fmt.Errorf("COMP_LINE environment variable not set")
	}

	// Split the line into words
	words := strings.Fields(line)
	if len(words) <= 1 {
		// If there's only one word (the command itself), output available flags
		fmt.Println("--version")
		fmt.Println("-v")
		fmt.Println("--completion")
		fmt.Println("--completion-output")
		return nil
	}

	// If there are more words, handle completion for specific flags
	lastWord := words[len(words)-1]
	if strings.HasPrefix("--completion", lastWord) {
		fmt.Println("--completion")
		return nil
	}
	if strings.HasPrefix("--completion-output", lastWord) {
		fmt.Println("--completion-output")
		return nil
	}
	if words[len(words)-2] == "--completion" {
		fmt.Println("bash")
		fmt.Println("zsh")
		fmt.Println("fish")
		return nil
	}

	return nil
}

func printVersion() {
	versionInfo := version.Map()

	// Print as JSON if stdout is not a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(versionInfo); err != nil {
			log.Fatalf("Error encoding version information: %v", err)
		}
		return
	}

	// Print in a human-readable format if stdout is a terminal
	fmt.Printf("Fair Project Server %s\n", versionInfo["version"])
	fmt.Printf("  Major: %d\n", versionInfo["major"])
	fmt.Printf("  Minor: %d\n", versionInfo["minor"])
	fmt.Printf("  Patch: %d\n", versionInfo["patch"])
	if hash, ok := versionInfo["branch_hash"].(string); ok && hash != "unknown" {
		fmt.Printf("  Branch Hash: %s\n", hash)
	}
}

func main() {
	// Define global flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	versionFlagShort := flag.Bool("v", false, "Print version information and exit (shorthand)")
	completionFlag := flag.String("completion", "", "Generate shell completion script (bash, zsh, fish)")
	completionOutputFlag := flag.String("completion-output", "", "Output file for completion script")
	flag.Parse()

	// Check if version flag is set
	if *versionFlag || *versionFlagShort {
		printVersion()
		return
	}

	// Check if completion flag is set
	if *completionFlag != "" {
		if err := handleCompletion(*completionFlag, *completionOutputFlag); err != nil {
			log.Fatalf("Error generating completion script: %v", err)
		}
		return
	}

	// Check if this is a completion request
	if len(os.Args) > 1 && os.Args[1] == "__complete" {
		if err := handleCompletionRequest(); err != nil {
			log.Fatalf("Error handling completion request: %v", err)
		}
		return
	}

	// 1. Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// 1.5. Load environment variables from .env file
	loadEnvFile()

	// 2. Initialize configuration from environment variables
	config.InitFromEnv()
	log.Printf("Log level set to: %s", config.GetLogLevel())

	// 3. Get the data directory from environment variable or use default
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	log.Printf("Using data directory: %s", dataDir)

	// 4. Create a storage manager with file-based storage
	ctx := context.Background()
	storageConfig := map[string]string{
		"baseDir": dataDir,
	}
	storageManager, err := storage.NewManager(ctx, storage.StorageTypeFile, storageConfig)
	if err != nil {
		log.Fatalf("Failed to create storage manager: %v", err)
	}
	defer storageManager.Close(ctx)

	// 5. Initialize default users
	if err := storageManager.InitializeDefaultUsers(ctx); err != nil {
		log.Fatalf("Failed to initialize default users: %v", err)
	}

	// 6. Set JWT secret from environment variable or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		authpkg.SetJWTSecret(jwtSecret)
	}

	// 7. Create an API handler with the storage manager
	apiHandler := api.NewHandler(storageManager)

	// 8. Create a handler with the logging and authentication middleware
	handler := logging.Middleware(authmiddleware.Middleware(http.HandlerFunc(apiHandler.MasterRouter)))

	// 9. Get the port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	// 10. Start the HTTP server with the middleware
	log.Printf("Starting server on port %s. Listening for requests on / , /classes, and /classes/...\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
