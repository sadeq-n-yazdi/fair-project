package errors

import (
	"fmt"
	"log"
	"net/http"
)

// ErrorCode is a unique identifier for an error type
type ErrorCode string

// Error codes for different types of errors
const (
	// Validation errors
	ErrInvalidClassName    ErrorCode = "INVALID_CLASS_NAME"
	ErrInvalidAssignmentID ErrorCode = "INVALID_ASSIGNMENT_ID"
	ErrInvalidRequestBody  ErrorCode = "INVALID_REQUEST_BODY"
	ErrEmptyClassName      ErrorCode = "EMPTY_CLASS_NAME"
	ErrMethodNotAllowed    ErrorCode = "METHOD_NOT_ALLOWED"

	// Not found errors
	ErrClassNotFound      ErrorCode = "CLASS_NOT_FOUND"
	ErrProjectsNotFound   ErrorCode = "PROJECTS_NOT_FOUND"
	ErrStudentsNotFound   ErrorCode = "STUDENTS_NOT_FOUND"
	ErrAssignmentNotFound ErrorCode = "ASSIGNMENT_NOT_FOUND"
	ErrResourceNotFound   ErrorCode = "RESOURCE_NOT_FOUND"

	// Data errors
	ErrNoProjects ErrorCode = "NO_PROJECTS"
	ErrNoStudents ErrorCode = "NO_STUDENTS"

	// Internal errors
	ErrInternal                ErrorCode = "INTERNAL_ERROR"
	ErrFailedToSaveProjects    ErrorCode = "FAILED_TO_SAVE_PROJECTS"
	ErrFailedToLoadProjects    ErrorCode = "FAILED_TO_LOAD_PROJECTS"
	ErrFailedToSaveStudents    ErrorCode = "FAILED_TO_SAVE_STUDENTS"
	ErrFailedToLoadStudents    ErrorCode = "FAILED_TO_LOAD_STUDENTS"
	ErrFailedToSaveAssignment  ErrorCode = "FAILED_TO_SAVE_ASSIGNMENT"
	ErrFailedToLoadAssignment  ErrorCode = "FAILED_TO_LOAD_ASSIGNMENT"
	ErrFailedToListAssignments ErrorCode = "FAILED_TO_LIST_ASSIGNMENTS"
)

// AppError represents an application error with a code and user-friendly message
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with the given code and message
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with a code and message
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// LogError logs the detailed error information
func LogError(err error) {
	if appErr, ok := err.(*AppError); ok {
		if appErr.Err != nil {
			log.Printf("[ERROR] %s: %s - %v", appErr.Code, appErr.Message, appErr.Err)
		} else {
			log.Printf("[ERROR] %s: %s", appErr.Code, appErr.Message)
		}
	} else {
		log.Printf("[ERROR] %v", err)
	}
}

// HTTPStatusFromErrorCode returns the appropriate HTTP status code for an error code
func HTTPStatusFromErrorCode(code ErrorCode) int {
	switch code {
	case ErrInvalidClassName, ErrInvalidAssignmentID, ErrInvalidRequestBody, ErrEmptyClassName:
		return http.StatusBadRequest
	case ErrClassNotFound, ErrProjectsNotFound, ErrStudentsNotFound, ErrAssignmentNotFound, ErrResourceNotFound:
		return http.StatusNotFound
	case ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case ErrNoProjects, ErrNoStudents:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// ToHTTPResponse converts an error to an HTTP response
func ToHTTPResponse(w http.ResponseWriter, err error) {
	var code ErrorCode
	var message string
	var status int

	if appErr, ok := err.(*AppError); ok {
		code = appErr.Code
		message = appErr.Message
		status = HTTPStatusFromErrorCode(code)
		// Log the detailed error
		LogError(appErr)
	} else {
		// If it's not an AppError, treat it as an internal error
		code = ErrInternal
		message = "An internal error occurred"
		status = http.StatusInternalServerError
		// Log the original error
		log.Printf("[ERROR] %v", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":{"code":"%s","message":"%s"}}`, code, message)
}
