package logging

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sadeq/fair-project-go/pkg/config"
)

func TestMiddleware(t *testing.T) {
	// Save the original log output and restore it after the test
	originalOutput := log.Writer()
	defer log.SetOutput(originalOutput)

	// Create a buffer to capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Create a test handler that returns a simple response
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello, World!"}`))
	})

	// Create a test server with the middleware
	server := httptest.NewServer(Middleware(testHandler))
	defer server.Close()

	// Test cases for different log levels
	testCases := []struct {
		name           string
		logLevel       config.LogLevel
		expectedLogs   []string
		unexpectedLogs []string
	}{
		{
			name:           "LogLevelNone",
			logLevel:       config.LogLevelNone,
			expectedLogs:   []string{},
			unexpectedLogs: []string{"[INFO]", "[DEBUG]"},
		},
		{
			name:           "LogLevelError",
			logLevel:       config.LogLevelError,
			expectedLogs:   []string{},
			unexpectedLogs: []string{"[INFO]", "[DEBUG]"},
		},
		{
			name:           "LogLevelInfo",
			logLevel:       config.LogLevelInfo,
			expectedLogs:   []string{"[INFO] Request:", "[INFO] Response:"},
			unexpectedLogs: []string{"[DEBUG]"},
		},
		{
			name:           "LogLevelDebug",
			logLevel:       config.LogLevelDebug,
			expectedLogs:   []string{"[INFO] Request:", "[INFO] Response:", "[DEBUG] Request Headers:", "[DEBUG] Response Body:"},
			unexpectedLogs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the log buffer
			logOutput.Reset()

			// Set the log level
			config.SetLogLevel(tc.logLevel)

			// Create a request with a body
			reqBody := bytes.NewBufferString(`{"test":"data"}`)
			req, err := http.NewRequest("POST", server.URL, reqBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Read the response body
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Check the logs
			logs := logOutput.String()
			for _, expected := range tc.expectedLogs {
				if !strings.Contains(logs, expected) {
					t.Errorf("Expected log to contain %q, but it didn't. Logs: %s", expected, logs)
				}
			}
			for _, unexpected := range tc.unexpectedLogs {
				if strings.Contains(logs, unexpected) {
					t.Errorf("Expected log not to contain %q, but it did. Logs: %s", unexpected, logs)
				}
			}
		})
	}
}

func TestMiddlewareFunc(t *testing.T) {
	// Save the original log output and restore it after the test
	originalOutput := log.Writer()
	defer log.SetOutput(originalOutput)

	// Create a buffer to capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Set the log level to debug
	config.SetLogLevel(config.LogLevelDebug)

	// Create a test handler that returns a simple response
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello, World!"}`))
	})

	// Create a test server with the middleware function
	server := httptest.NewServer(MiddlewareFunc(testHandler))
	defer server.Close()

	// Create a request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the logs
	logs := logOutput.String()
	expectedLogs := []string{"[INFO] Request:", "[INFO] Response:", "[DEBUG] Response Body:"}
	for _, expected := range expectedLogs {
		if !strings.Contains(logs, expected) {
			t.Errorf("Expected log to contain %q, but it didn't. Logs: %s", expected, logs)
		}
	}
}
