package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// Factory implements the repository.RepositoryFactory interface for SQLite storage
type Factory struct {
	dbPath string
	db     *sql.DB
	mu     sync.Mutex

	userRepo       repository.UserRepository
	classRepo      repository.ClassRepository
	projectRepo    repository.ProjectRepository
	studentRepo    repository.StudentRepository
	assignmentRepo repository.AssignmentRepository
}

// NewFactory creates a new SQLite repository factory
func NewFactory(dbPath string) *Factory {
	return &Factory{
		dbPath: dbPath,
	}
}

// initialize initializes the SQLite database connection
func (f *Factory) initialize(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.db != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", f.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	f.db = db
	return nil
}

// Close closes the SQLite database connection
func (f *Factory) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.db != nil {
		err := f.db.Close()
		f.db = nil
		return err
	}
	return nil
}

// CreateUserRepository creates a SQLite user repository
func (f *Factory) CreateUserRepository(ctx context.Context) (repository.UserRepository, error) {
	if err := f.initialize(ctx); err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db)
		if err := f.userRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.userRepo, nil
}

// CreateClassRepository creates a SQLite class repository
func (f *Factory) CreateClassRepository(ctx context.Context) (repository.ClassRepository, error) {
	if err := f.initialize(ctx); err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.classRepo == nil {
		f.classRepo = NewClassRepository(f.db)
		if err := f.classRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.classRepo, nil
}

// CreateProjectRepository creates a SQLite project repository
func (f *Factory) CreateProjectRepository(ctx context.Context) (repository.ProjectRepository, error) {
	if err := f.initialize(ctx); err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.projectRepo == nil {
		f.projectRepo = NewProjectRepository(f.db)
		if err := f.projectRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.projectRepo, nil
}

// CreateStudentRepository creates a SQLite student repository
func (f *Factory) CreateStudentRepository(ctx context.Context) (repository.StudentRepository, error) {
	if err := f.initialize(ctx); err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.studentRepo == nil {
		f.studentRepo = NewStudentRepository(f.db)
		if err := f.studentRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.studentRepo, nil
}

// CreateAssignmentRepository creates a SQLite assignment repository
func (f *Factory) CreateAssignmentRepository(ctx context.Context) (repository.AssignmentRepository, error) {
	if err := f.initialize(ctx); err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.assignmentRepo == nil {
		f.assignmentRepo = NewAssignmentRepository(f.db)
		if err := f.assignmentRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.assignmentRepo, nil
}
