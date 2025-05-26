package assignment

import (
	"math/rand"
	_ "time"

	"github.com/sadeq/fair-project-go/pkg/models"
)

// assignProjects assigns students to projects based on their preferences.
// It attempts to give each student their first available choice.
// The process is iterative and involves shuffling students to ensure fairness.
func AssignProjects(students []models.Student, projects []string) (map[string]string, map[string]string, []models.Student, []string) {
	// Seed the random number generator once.
	// Note: In a larger application, this should ideally be done once in the main function.
	// For this specific function structure, and to ensure it's seeded if called independently,
	// it's placed here. If main.go calls this, the seed in main.go will suffice.
	// rand.Seed(time.Now().UnixNano()) // Moved to main.go as per later instruction

	selectedByProjects := make(map[string]string)
	selectedByStudent := make(map[string]string)

	// Create a map for quick lookups of available projects
	availableProjects := make(map[string]bool)
	for _, p := range projects {
		availableProjects[p] = true
	}

	maxIterations := len(projects)
	if len(students) > maxIterations {
		maxIterations = len(students)
	}

	for i := 0; i < maxIterations; i++ {
		// Shuffle students for fairness in each iteration
		rand.Shuffle(len(students), func(i, j int) {
			students[i], students[j] = students[j], students[i]
		})

		allAssignedOrNoPreferences := true
		for studentIdx := range students {
			student := &students[studentIdx] // Use a pointer to modify the original student struct

			if student.Assigned {
				continue
			}
			if len(student.Preferences) == 0 {
				continue
			}

			allAssignedOrNoPreferences = false // Found a student who is not assigned and has preferences

			firstChoice := student.Preferences[0]

			if availableProjects[firstChoice] {
				selectedByProjects[firstChoice] = student.Name
				selectedByStudent[student.Name] = firstChoice
				student.Assigned = true
				student.AssignedProject = firstChoice

				delete(availableProjects, firstChoice) // Remove project from available map

			} else {
				// Add to bad preferences and remove from current preferences
				student.BadPreferences = append(student.BadPreferences, firstChoice)
				student.Preferences = student.Preferences[1:]
			}
		}
		if allAssignedOrNoPreferences {
			// Optimization: if all students are assigned or have no preferences left, no need to continue iterations.
			break
		}
	}

	// Update the list of remaining projects based on the availableProjects map
	var remainingProjects []string
	for proj := range availableProjects {
		remainingProjects = append(remainingProjects, proj)
	}
	// Note: The 'projects' slice was already modified in place when a project was assigned.
	// 'remainingProjects' from availableProjects map is a more reliable way to get the final list.

	return selectedByProjects, selectedByStudent, students, remainingProjects
}

// Utility function to seed random number generator, if not already done in main.
// This is a safeguard. Ideally, main should call Seed.
func EnsureRandSeeded() {
	// This is a common pattern, but time.Now().UnixNano() might not be unique enough
	// if called in rapid succession. For this assignment, it's generally fine.
	// A global sync.Once could also be used in a larger app to ensure seeding happens exactly once.
	// For now, assuming main will handle it or this will be sufficient.
	// rand.Seed(time.Now().UnixNano()) // Intentionally commented out, main.go will handle this.
}

func init() {
	// It's common to seed in init(), but for testing or specific scenarios,
	// explicit seeding (e.g., in main or at the start of assignProjects if it were standalone)
	// is often preferred for better control.
	// rand.Seed(time.Now().UnixNano()) // Moving seed to main for better control
}
