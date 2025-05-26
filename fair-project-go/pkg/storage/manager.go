package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
	"github.com/sadeq/fair-project-go/pkg/storage/repository/file"
)

// StorageType represents the type of storage to use
type StorageType string

const (
	// StorageTypeFile represents file-based storage
	StorageTypeFile StorageType = "file"
	// StorageTypeSQLite represents SQLite storage
	StorageTypeSQLite StorageType = "sqlite"
	// StorageTypePostgres represents PostgreSQL storage
	StorageTypePostgres StorageType = "postgres"
)

// Manager manages the storage repositories
type Manager struct {
	factory repository.RepositoryFactory
	mu      sync.Mutex

	userRepo       repository.UserRepository
	classRepo      repository.ClassRepository
	projectRepo    repository.ProjectRepository
	studentRepo    repository.StudentRepository
	assignmentRepo repository.AssignmentRepository
}

// NewManager creates a new storage manager
func NewManager(ctx context.Context, storageType StorageType, config map[string]string) (*Manager, error) {
	var factory repository.RepositoryFactory
	var err error

	switch storageType {
	case StorageTypeFile:
		baseDir := config["baseDir"]
		if baseDir == "" {
			baseDir = "data"
		}
		factory = file.NewFactory(baseDir)
	case StorageTypeSQLite:
		// SQLite factory would be created here
		return nil, fmt.Errorf("SQLite storage not implemented yet")
	case StorageTypePostgres:
		// PostgreSQL factory would be created here
		return nil, fmt.Errorf("PostgreSQL storage not implemented yet")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}

	manager := &Manager{
		factory: factory,
	}

	// Initialize the repositories
	if err = manager.initRepositories(ctx); err != nil {
		return nil, err
	}

	return manager, nil
}

// initRepositories initializes all repositories
func (m *Manager) initRepositories(ctx context.Context) error {
	var err error

	// Initialize user repository
	m.userRepo, err = m.factory.CreateUserRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user repository: %w", err)
	}

	// Initialize class repository
	m.classRepo, err = m.factory.CreateClassRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class repository: %w", err)
	}

	// Initialize project repository
	m.projectRepo, err = m.factory.CreateProjectRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to create project repository: %w", err)
	}

	// Initialize student repository
	m.studentRepo, err = m.factory.CreateStudentRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to create student repository: %w", err)
	}

	// Initialize assignment repository
	m.assignmentRepo, err = m.factory.CreateAssignmentRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to create assignment repository: %w", err)
	}

	return nil
}

// Close closes all repositories
func (m *Manager) Close(ctx context.Context) error {
	var errs []error

	if m.userRepo != nil {
		if err := m.userRepo.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close user repository: %w", err))
		}
	}

	if m.classRepo != nil {
		if err := m.classRepo.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close class repository: %w", err))
		}
	}

	if m.projectRepo != nil {
		if err := m.projectRepo.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close project repository: %w", err))
		}
	}

	if m.studentRepo != nil {
		if err := m.studentRepo.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close student repository: %w", err))
		}
	}

	if m.assignmentRepo != nil {
		if err := m.assignmentRepo.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close assignment repository: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing repositories: %v", errs)
	}

	return nil
}

// UserRepository returns the user repository
func (m *Manager) UserRepository() repository.UserRepository {
	return m.userRepo
}

// ClassRepository returns the class repository
func (m *Manager) ClassRepository() repository.ClassRepository {
	return m.classRepo
}

// ProjectRepository returns the project repository
func (m *Manager) ProjectRepository() repository.ProjectRepository {
	return m.projectRepo
}

// StudentRepository returns the student repository
func (m *Manager) StudentRepository() repository.StudentRepository {
	return m.studentRepo
}

// AssignmentRepository returns the assignment repository
func (m *Manager) AssignmentRepository() repository.AssignmentRepository {
	return m.assignmentRepo
}

// InitializeDefaultUsers initializes default users if no users exist
func (m *Manager) InitializeDefaultUsers(ctx context.Context) error {
	return m.userRepo.InitializeDefaultUsers(ctx)
}
