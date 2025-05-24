package main

// Project is a type alias for string, representing the name of a project.
type Project string

// Student represents a student with their preferences and assignment status.
type Student struct {
	Name            string
	Preferences     []string
	Assigned        bool
	AssignedProject string
	BadPreferences  []string
}
