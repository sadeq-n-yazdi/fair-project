package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sadeq/fair-project-go/pkg/cli"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Ensure the base data directory exists
	if err := storage.EnsureBaseDir(); err != nil {
		log.Fatalf("Failed to ensure base data directory: %v", err)
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

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  fair-project-cli <command> [options]")
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
	fmt.Println("  help               Show this help message")
	fmt.Println("\nRun 'fair-project-cli <command> -h' for more information on a command.")
}
