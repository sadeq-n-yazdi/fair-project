package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/sadeq/fair-project-go/pkg/auth"
	"github.com/sadeq/fair-project-go/pkg/errors"
	"github.com/sadeq/fair-project-go/pkg/models"
	"github.com/sadeq/fair-project-go/pkg/storage"
)

// contextKey is a type for context keys
type contextKey string

// UserKey is the key for the user in the request context
const UserKey contextKey = "user"

// Add authentication error codes
const (
	ErrUnauthorized      errors.ErrorCode = "UNAUTHORIZED"
	ErrForbidden         errors.ErrorCode = "FORBIDDEN"
	ErrInvalidToken      errors.ErrorCode = "INVALID_TOKEN"
	ErrMissingToken      errors.ErrorCode = "MISSING_TOKEN"
	ErrTokenExpired      errors.ErrorCode = "TOKEN_EXPIRED"
	ErrInvalidAuthHeader errors.ErrorCode = "INVALID_AUTH_HEADER"
)

// GetUserFromContext gets the user from the request context
func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserKey).(*models.User)
	return user, ok
}

// Middleware returns a middleware function that validates JWT tokens
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for docs, root, and specific endpoints
		if r.URL.Path == "/" || r.URL.Path == "/docs" || r.URL.Path == "/version" ||
			strings.HasPrefix(r.URL.Path, "/auth") || r.URL.Path == "/user/whoami" {
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.ToHTTPResponse(w, errors.New(ErrMissingToken, "Missing authorization token"))
			return
		}

		// Check if the Authorization header has the correct format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			errors.ToHTTPResponse(w, errors.New(ErrInvalidAuthHeader, "Invalid authorization header format"))
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			errors.ToHTTPResponse(w, errors.Wrap(err, ErrInvalidToken, "Invalid token"))
			return
		}

		// Get the user from the database
		username := claims.Username
		user, exists, err := storage.GetUser(username)
		if err != nil {
			errors.ToHTTPResponse(w, errors.Wrap(err, ErrUnauthorized, "Failed to get user"))
			return
		}

		if !exists {
			errors.ToHTTPResponse(w, errors.New(ErrUnauthorized, "User not found"))
			return
		}

		if !user.Enabled {
			errors.ToHTTPResponse(w, errors.New(ErrUnauthorized, "User is disabled"))
			return
		}

		// Convert to pointer for context
		userPtr := &user

		// Add the user to the request context
		ctx := context.WithValue(r.Context(), UserKey, userPtr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoles returns a middleware function that requires the user to have specific roles
func RequireRoles(roles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r)
			if !ok {
				errors.ToHTTPResponse(w, errors.New(ErrUnauthorized, "User not found in context"))
				return
			}

			if !user.HasAnyRole(roles...) {
				errors.ToHTTPResponse(w, errors.New(ErrForbidden, "Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin returns a middleware function that requires the user to be an admin or superadmin
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRoles(models.RoleAdmin, models.RoleSuperAdmin)(next)
}

// RequireSuperAdmin returns a middleware function that requires the user to be a superadmin
func RequireSuperAdmin(next http.Handler) http.Handler {
	return RequireRoles(models.RoleSuperAdmin)(next)
}

// RequireTeacher returns a middleware function that requires the user to be a teacher, admin, or superadmin
func RequireTeacher(next http.Handler) http.Handler {
	return RequireRoles(models.RoleTeacher, models.RoleAdmin, models.RoleSuperAdmin)(next)
}

// RequireStudent returns a middleware function that requires the user to be a student, teacher, admin, or superadmin
func RequireStudent(next http.Handler) http.Handler {
	return RequireRoles(models.RoleStudent, models.RoleTeacher, models.RoleAdmin, models.RoleSuperAdmin)(next)
}
