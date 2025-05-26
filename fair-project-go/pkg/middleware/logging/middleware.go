package logging

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sadeq/fair-project-go/pkg/config"
)

// responseWriterWrapper is a wrapper for http.ResponseWriter that captures the status code and response body
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	buffer     *bytes.Buffer
}

// newResponseWriterWrapper creates a new responseWriterWrapper
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
		buffer:         &bytes.Buffer{},
	}
}

// WriteHeader captures the status code and calls the underlying ResponseWriter's WriteHeader
func (rww *responseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body and calls the underlying ResponseWriter's Write
func (rww *responseWriterWrapper) Write(b []byte) (int, error) {
	// Write to the buffer for logging
	rww.buffer.Write(b)
	// Write to the original response writer
	return rww.ResponseWriter.Write(b)
}

// GetStatusCode returns the captured status code
func (rww *responseWriterWrapper) GetStatusCode() int {
	return rww.statusCode
}

// GetResponseBody returns the captured response body
func (rww *responseWriterWrapper) GetResponseBody() string {
	return rww.buffer.String()
}

// Middleware returns a middleware function that logs HTTP requests and responses
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging if log level is none
		if config.GetLogLevel() == config.LogLevelNone {
			next.ServeHTTP(w, r)
			return
		}

		// Start time
		startTime := time.Now()

		// Log request details
		if config.ShouldLog(config.LogLevelInfo) {
			log.Printf("[INFO] Request: %s %s", r.Method, r.URL.Path)
		}

		// Log request headers if debug level
		if config.ShouldLog(config.LogLevelDebug) {
			log.Printf("[DEBUG] Request Headers: %v", r.Header)
		}

		// Log request body if debug level and body exists
		if config.ShouldLog(config.LogLevelDebug) && r.Body != nil {
			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("[ERROR] Failed to read request body: %v", err)
			} else if len(bodyBytes) > 0 {
				// Log the body
				log.Printf("[DEBUG] Request Body: %s", string(bodyBytes))
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create a response writer wrapper to capture the response
		rww := newResponseWriterWrapper(w)

		// Call the next handler
		next.ServeHTTP(rww, r)

		// Calculate duration
		duration := time.Since(startTime)

		// Log response details
		if config.ShouldLog(config.LogLevelInfo) {
			log.Printf("[INFO] Response: %s %s - Status: %d - Duration: %v",
				r.Method, r.URL.Path, rww.GetStatusCode(), duration)
		}

		// Log response body if debug level
		if config.ShouldLog(config.LogLevelDebug) {
			log.Printf("[DEBUG] Response Body: %s", rww.GetResponseBody())
		}
	})
}

// MiddlewareFunc is a convenience function that returns a middleware function
// that can be used with http.HandleFunc
func MiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Middleware(next).ServeHTTP(w, r)
	}
}
