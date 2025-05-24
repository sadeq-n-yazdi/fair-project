package main

import "fmt"

// printResults prints the assignment results to the console.
func printResults(
	selectedByProjects map[string]string,
	selectedByStudent map[string]string,
	students []Student,
	remainingProjects []string,
) {
	fmt.Println("========================")

	fmt.Println("Selected by Projects:")
	if len(selectedByProjects) == 0 {
		fmt.Println("  (No projects were assigned)")
	} else {
		for project, student := range selectedByProjects {
			fmt.Printf("  %s: %s\n", project, student)
		}
	}

	fmt.Println("------------------------")

	fmt.Println("Selected by Student:")
	if len(selectedByStudent) == 0 {
		fmt.Println("  (No students were assigned)")
	} else {
		for student, project := range selectedByStudent {
			fmt.Printf("  %s: %s\n", student, project)
		}
	}

	fmt.Println("------------------------")

	fmt.Println("Remaining Projects:")
	if len(remainingProjects) == 0 {
		fmt.Println("  (All projects were assigned)")
	} else {
		for _, project := range remainingProjects {
			fmt.Printf("  - %s\n", project)
		}
	}

	fmt.Println("------------------------")

	fmt.Println("Student Details:")
	if len(students) == 0 {
		fmt.Println("  (No students to display)")
	} else {
		for _, student := range students {
			fmt.Printf("  Student: %s\n", student.Name)
			if student.Assigned {
				fmt.Printf("    Assigned Project: %s\n", student.AssignedProject)
			} else {
				fmt.Println("    Status: Not Assigned")
				if len(student.Preferences) > 0 {
					fmt.Printf("    Remaining Preferences: %v\n", student.Preferences)
				} else {
					fmt.Println("    Remaining Preferences: None")
				}
				if len(student.BadPreferences) > 0 {
					fmt.Printf("    Bad Preferences (tried and failed): %v\n", student.BadPreferences)
				}
			}
		}
	}
	fmt.Println("------------------------")
}
