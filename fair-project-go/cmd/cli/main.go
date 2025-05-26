package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/sadeq/fair-project-go/pkg/cli"
	"github.com/sadeq/fair-project-go/pkg/storage"
	"log"
	"os"
	"strings"
)

func main() {

	// Ensure the base data directory exists
	if err := storage.EnsureBaseDir(); err != nil {
		log.Fatalf("Failed to ensure base data directory: %v", err)
	}

	// Load environment variables from .env file
	loadEnvFile()

	// Initialize default users if none exist
	if err := storage.InitializeDefaultUsers(); err != nil {
		log.Fatalf("Failed to initialize default users: %v", err)
	}

	// Prompt for superadmin creation if none exists
	if err := cli.PromptForSuperAdmin(); err != nil {
		log.Printf("Warning: Failed to prompt for superadmin: %v", err)
	}

	// Define command-line flags
	createClassCmd := flag.NewFlagSet("create-class", flag.ExitOnError)
	createClassName := createClassCmd.String("name", "", "Name of the class/term to create")

	listClassesCmd := flag.NewFlagSet("list-classes", flag.ExitOnError)

	importProjectsCmd := flag.NewFlagSet("import-projects", flag.ExitOnError)
	importProjectsClass := importProjectsCmd.String("class", "", "Name of the class/term")
	importProjectsFile := importProjectsCmd.String("file", "", "Path to the projects file")

	importStudentsCmd := flag.NewFlagSet("import-students", flag.ExitOnError)
	importStudentsClass := importStudentsCmd.String("class", "", "Name of the class/term")
	importStudentsFile := importStudentsCmd.String("file", "", "Path to the students file")

	listProjectsCmd := flag.NewFlagSet("list-projects", flag.ExitOnError)
	listProjectsClass := listProjectsCmd.String("class", "", "Name of the class/term")

	listStudentsCmd := flag.NewFlagSet("list-students", flag.ExitOnError)
	listStudentsClass := listStudentsCmd.String("class", "", "Name of the class/term")

	runAssignmentCmd := flag.NewFlagSet("run-assignment", flag.ExitOnError)
	runAssignmentClass := runAssignmentCmd.String("class", "", "Name of the class/term")

	listAssignmentsCmd := flag.NewFlagSet("list-assignments", flag.ExitOnError)
	listAssignmentsClass := listAssignmentsCmd.String("class", "", "Name of the class/term")

	showAssignmentCmd := flag.NewFlagSet("show-assignment", flag.ExitOnError)
	showAssignmentClass := showAssignmentCmd.String("class", "", "Name of the class/term")
	showAssignmentID := showAssignmentCmd.String("id", "", "ID of the assignment")

	exportAssignmentCmd := flag.NewFlagSet("export-assignment", flag.ExitOnError)
	exportAssignmentClass := exportAssignmentCmd.String("class", "", "Name of the class/term")
	exportAssignmentID := exportAssignmentCmd.String("id", "", "ID of the assignment")
	exportAssignmentOutput := exportAssignmentCmd.String("output", "", "Path to the output file")

	// Check if a command was provided
	if len(os.Args) < 2 {
		fmt.Println("Expected a command")
		printUsage()
		os.Exit(1)
	}

	// Define user management command flags
	createSuperAdminCmd := flag.NewFlagSet("create-superadmin", flag.ExitOnError)
	createSuperAdminUsername := createSuperAdminCmd.String("username", "", "Username for the superadmin")
	createSuperAdminPassword := createSuperAdminCmd.String("password", "", "Password for the superadmin")
	createSuperAdminKey := createSuperAdminCmd.String("key", "", "Superadmin key from .env file")

	changePasswordCmd := flag.NewFlagSet("change-password", flag.ExitOnError)
	changePasswordUsername := changePasswordCmd.String("username", "", "Username of the user")
	changePasswordNewPassword := changePasswordCmd.String("password", "", "New password for the user")

	listUsersCmd := flag.NewFlagSet("list-users", flag.ExitOnError)

	// Parse the command
	switch os.Args[1] {
	case "create-class":
		createClassCmd.Parse(os.Args[2:])
		if *createClassName == "" {
			fmt.Println("Error: class name is required")
			createClassCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.CreateClassTerm(*createClassName)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list-classes":
		listClassesCmd.Parse(os.Args[2:])
		err := cli.ListClassTerms()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "import-projects":
		importProjectsCmd.Parse(os.Args[2:])
		if *importProjectsClass == "" {
			fmt.Println("Error: class name is required")
			importProjectsCmd.PrintDefaults()
			os.Exit(1)
		}
		if *importProjectsFile == "" {
			fmt.Println("Error: file path is required")
			importProjectsCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ImportProjects(*importProjectsClass, *importProjectsFile)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "import-students":
		importStudentsCmd.Parse(os.Args[2:])
		if *importStudentsClass == "" {
			fmt.Println("Error: class name is required")
			importStudentsCmd.PrintDefaults()
			os.Exit(1)
		}
		if *importStudentsFile == "" {
			fmt.Println("Error: file path is required")
			importStudentsCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ImportStudents(*importStudentsClass, *importStudentsFile)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list-projects":
		listProjectsCmd.Parse(os.Args[2:])
		if *listProjectsClass == "" {
			fmt.Println("Error: class name is required")
			listProjectsCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ListProjects(*listProjectsClass)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list-students":
		listStudentsCmd.Parse(os.Args[2:])
		if *listStudentsClass == "" {
			fmt.Println("Error: class name is required")
			listStudentsCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ListStudents(*listStudentsClass)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "run-assignment":
		runAssignmentCmd.Parse(os.Args[2:])
		if *runAssignmentClass == "" {
			fmt.Println("Error: class name is required")
			runAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.RunAssignment(*runAssignmentClass)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list-assignments":
		listAssignmentsCmd.Parse(os.Args[2:])
		if *listAssignmentsClass == "" {
			fmt.Println("Error: class name is required")
			listAssignmentsCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ListAssignments(*listAssignmentsClass)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "show-assignment":
		showAssignmentCmd.Parse(os.Args[2:])
		if *showAssignmentClass == "" {
			fmt.Println("Error: class name is required")
			showAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		if *showAssignmentID == "" {
			fmt.Println("Error: assignment ID is required")
			showAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ShowAssignment(*showAssignmentClass, *showAssignmentID)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "export-assignment":
		exportAssignmentCmd.Parse(os.Args[2:])
		if *exportAssignmentClass == "" {
			fmt.Println("Error: class name is required")
			exportAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		if *exportAssignmentID == "" {
			fmt.Println("Error: assignment ID is required")
			exportAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		if *exportAssignmentOutput == "" {
			fmt.Println("Error: output path is required")
			exportAssignmentCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ExportAssignment(*exportAssignmentClass, *exportAssignmentID, *exportAssignmentOutput)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "help":
		printUsage()

	case "create-superadmin":
		createSuperAdminCmd.Parse(os.Args[2:])
		if *createSuperAdminUsername == "" {
			fmt.Println("Error: username is required")
			createSuperAdminCmd.PrintDefaults()
			os.Exit(1)
		}
		if *createSuperAdminPassword == "" {
			fmt.Println("Error: password is required")
			createSuperAdminCmd.PrintDefaults()
			os.Exit(1)
		}
		if *createSuperAdminKey == "" {
			fmt.Println("Error: superadmin key is required")
			createSuperAdminCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.CreateSuperAdmin(*createSuperAdminUsername, *createSuperAdminPassword, *createSuperAdminKey)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "change-password":
		changePasswordCmd.Parse(os.Args[2:])
		if *changePasswordUsername == "" {
			fmt.Println("Error: username is required")
			changePasswordCmd.PrintDefaults()
			os.Exit(1)
		}
		if *changePasswordNewPassword == "" {
			fmt.Println("Error: new password is required")
			changePasswordCmd.PrintDefaults()
			os.Exit(1)
		}
		err := cli.ChangeUserPassword(*changePasswordUsername, *changePasswordNewPassword)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list-users":
		listUsersCmd.Parse(os.Args[2:])
		err := cli.ListUsers()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
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

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  fair-ctl <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  create-class       Create a new class/term")
	fmt.Println("  list-classes       List all class/terms")
	fmt.Println("  import-projects    Import projects from a file")
	fmt.Println("  import-students    Import students from a file")
	fmt.Println("  list-projects      List all projects for a class/term")
	fmt.Println("  list-students      List all students for a class/term")
	fmt.Println("  run-assignment     Run the assignment algorithm for a class/term")
	fmt.Println("  list-assignments   List all assignments for a class/term")
	fmt.Println("  show-assignment    Show the details of an assignment")
	fmt.Println("  export-assignment  Export an assignment to a JSON file")
	fmt.Println("  create-superadmin  Create the first superadmin user")
	fmt.Println("  change-password    Change a user's password")
	fmt.Println("  list-users         List all users")
	fmt.Println("  help               Show this help message")
	fmt.Println("\nRun 'fair-ctl <command> -h' for more information on a command.")
}
