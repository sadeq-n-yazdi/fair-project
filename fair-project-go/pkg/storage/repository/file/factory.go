package file

import (
	"context"
	"sync"

	"github.com/sadeq/fair-project-go/pkg/storage/repository"
)

// Factory implements the repository.RepositoryFactory interface for file-based storage
type Factory struct {
	baseDir string
	mu      sync.Mutex

	userRepo       repository.UserRepository
	classRepo      repository.ClassRepository
	projectRepo    repository.ProjectRepository
	studentRepo    repository.StudentRepository
	assignmentRepo repository.AssignmentRepository
}

// NewFactory creates a new file-based repository factory
func NewFactory(baseDir string) *Factory {
	return &Factory{
		baseDir: baseDir,
	}
}

// CreateUserRepository creates a file-based user repository
func (f *Factory) CreateUserRepository(ctx context.Context) (repository.UserRepository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.baseDir)
		if err := f.userRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.userRepo, nil
}

// CreateClassRepository creates a file-based class repository
func (f *Factory) CreateClassRepository(ctx context.Context) (repository.ClassRepository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.classRepo == nil {
		f.classRepo = NewClassRepository(f.baseDir)
		if err := f.classRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.classRepo, nil
}

// CreateProjectRepository creates a file-based project repository
func (f *Factory) CreateProjectRepository(ctx context.Context) (repository.ProjectRepository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.projectRepo == nil {
		f.projectRepo = NewProjectRepository(f.baseDir)
		if err := f.projectRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.projectRepo, nil
}

// CreateStudentRepository creates a file-based student repository
func (f *Factory) CreateStudentRepository(ctx context.Context) (repository.StudentRepository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.studentRepo == nil {
		f.studentRepo = NewStudentRepository(f.baseDir)
		if err := f.studentRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.studentRepo, nil
}

// CreateAssignmentRepository creates a file-based assignment repository
func (f *Factory) CreateAssignmentRepository(ctx context.Context) (repository.AssignmentRepository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.assignmentRepo == nil {
		f.assignmentRepo = NewAssignmentRepository(f.baseDir)
		if err := f.assignmentRepo.Initialize(ctx); err != nil {
			return nil, err
		}
	}
	return f.assignmentRepo, nil
}
