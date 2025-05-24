package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	// 1. Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// 2. Define filenames
	projectsFile := "projects.csv"
	studentsFile := "students.csv"

	// 3. Read projects
	projects, err := readProjectsFromFile(projectsFile)
	if err != nil {
		log.Fatalf("Error reading projects from file %s: %v", projectsFile, err)
	}
	if len(projects) == 0 {
		log.Fatalf("No projects found in %s. Exiting.", projectsFile)
	}
	fmt.Printf("Successfully read %d projects.\n", len(projects))

	// 4. Read students
	students, err := readStudentsFromFile(studentsFile)
	if err != nil {
		log.Fatalf("Error reading students from file %s: %v", studentsFile, err)
	}
	if len(students) == 0 {
		log.Fatalf("No students found in %s. Exiting.", studentsFile)
	}
	fmt.Printf("Successfully read %d students.\n", len(students))

	// 5. Perform assignments
	// Ensure students are passed as a slice of Student structs, and projects as a slice of strings.
	// The assignProjects function is expected to handle the logic correctly.
	selectedByProjects, selectedByStudent, updatedStudents, remainingProjects := assignProjects(students, projects)

	// 6. Print results
	printResults(selectedByProjects, selectedByStudent, updatedStudents, remainingProjects)

	fmt.Println("\nProgram finished.")
}
