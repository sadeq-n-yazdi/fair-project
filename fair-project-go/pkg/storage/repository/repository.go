package repository

import (
	"context"

	"github.com/sadeq/fair-project-go/pkg/models"
)

// Repository is the base interface for all repositories
type Repository interface {
	// Initialize initializes the repository
	Initialize(ctx context.Context) error
	// Close closes the repository
	Close(ctx context.Context) error
}

// UserRepository defines operations for user data
type UserRepository interface {
	Repository
	// GetUser gets a user by username
	GetUser(ctx context.Context, username string) (models.User, bool, error)
	// SaveUser saves a user
	SaveUser(ctx context.Context, user models.User) error
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, username string) error
	// GetAllUsers gets all users
	GetAllUsers(ctx context.Context) (map[string]models.User, error)
	// InitializeDefaultUsers creates default users if no users exist
	InitializeDefaultUsers(ctx context.Context) error
}

// ClassRepository defines operations for class/term data
type ClassRepository interface {
	Repository
	// CreateClassTerm creates a new class/term
	CreateClassTerm(ctx context.Context, className string) error
	// ListClassTerms lists all class/terms
	ListClassTerms(ctx context.Context) ([]string, error)
}

// ProjectRepository defines operations for project data
type ProjectRepository interface {
	Repository
	// SaveProjects saves projects for a class/term
	SaveProjects(ctx context.Context, className string, projects []string) error
	// LoadProjects loads projects for a class/term
	LoadProjects(ctx context.Context, className string) ([]string, error)
	// ReadProjectsFromFile reads projects from a file
	ReadProjectsFromFile(ctx context.Context, fileName string) ([]string, error)
}

// StudentRepository defines operations for student data
type StudentRepository interface {
	Repository
	// SaveStudents saves students for a class/term
	SaveStudents(ctx context.Context, className string, students []models.StudentInput) error
	// LoadStudents loads students for a class/term
	LoadStudents(ctx context.Context, className string) ([]models.StudentInput, error)
	// ReadStudentsFromFile reads students from a file
	ReadStudentsFromFile(ctx context.Context, fileName string) ([]models.Student, error)
}

// AssignmentRepository defines operations for assignment data
type AssignmentRepository interface {
	Repository
	// SaveAssignmentResults saves assignment results for a class/term
	SaveAssignmentResults(ctx context.Context, className string, assignmentID string, results models.AssignmentOutput) error
	// ListAssignmentResults lists all assignment results for a class/term
	ListAssignmentResults(ctx context.Context, className string) ([]string, error)
	// LoadAssignmentResult loads an assignment result for a class/term
	LoadAssignmentResult(ctx context.Context, className string, assignmentID string) (models.AssignmentOutput, error)
}

// RepositoryFactory creates repositories for different storage types
type RepositoryFactory interface {
	// CreateUserRepository creates a user repository
	CreateUserRepository(ctx context.Context) (UserRepository, error)
	// CreateClassRepository creates a class repository
	CreateClassRepository(ctx context.Context) (ClassRepository, error)
	// CreateProjectRepository creates a project repository
	CreateProjectRepository(ctx context.Context) (ProjectRepository, error)
	// CreateStudentRepository creates a student repository
	CreateStudentRepository(ctx context.Context) (StudentRepository, error)
	// CreateAssignmentRepository creates an assignment repository
	CreateAssignmentRepository(ctx context.Context) (AssignmentRepository, error)
}
