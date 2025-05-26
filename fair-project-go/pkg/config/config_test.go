package config

import (
	"os"
	"testing"
)

func TestLogLevelFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected LogLevel
	}{
		{"none", LogLevelNone},
		{"error", LogLevelError},
		{"info", LogLevelInfo},
		{"debug", LogLevelDebug},
		{"NONE", LogLevelNone},
		{"ERROR", LogLevelError},
		{"INFO", LogLevelInfo},
		{"DEBUG", LogLevelDebug},
		{"invalid", LogLevelInfo}, // Default to info for invalid input
		{"", LogLevelInfo},        // Default to info for empty input
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := LogLevelFromString(tc.input)
			if result != tc.expected {
				t.Errorf("LogLevelFromString(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSetAndGetLogLevel(t *testing.T) {
	// Test setting and getting log level
	testCases := []LogLevel{
		LogLevelNone,
		LogLevelError,
		LogLevelInfo,
		LogLevelDebug,
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(t *testing.T) {
			SetLogLevel(tc)
			result := GetLogLevel()
			if result != tc {
				t.Errorf("GetLogLevel() = %v, expected %v", result, tc)
			}
		})
	}
}

func TestShouldLog(t *testing.T) {
	testCases := []struct {
		currentLevel LogLevel
		checkLevel   LogLevel
		expected     bool
	}{
		{LogLevelNone, LogLevelNone, true},
		{LogLevelNone, LogLevelError, false},
		{LogLevelNone, LogLevelInfo, false},
		{LogLevelNone, LogLevelDebug, false},

		{LogLevelError, LogLevelNone, true},
		{LogLevelError, LogLevelError, true},
		{LogLevelError, LogLevelInfo, false},
		{LogLevelError, LogLevelDebug, false},

		{LogLevelInfo, LogLevelNone, true},
		{LogLevelInfo, LogLevelError, true},
		{LogLevelInfo, LogLevelInfo, true},
		{LogLevelInfo, LogLevelDebug, false},

		{LogLevelDebug, LogLevelNone, true},
		{LogLevelDebug, LogLevelError, true},
		{LogLevelDebug, LogLevelInfo, true},
		{LogLevelDebug, LogLevelDebug, true},
	}

	for _, tc := range testCases {
		t.Run(tc.currentLevel.String()+"_"+tc.checkLevel.String(), func(t *testing.T) {
			SetLogLevel(tc.currentLevel)
			result := ShouldLog(tc.checkLevel)
			if result != tc.expected {
				t.Errorf("ShouldLog(%v) with current level %v = %v, expected %v", tc.checkLevel, tc.currentLevel, result, tc.expected)
			}
		})
	}
}

func TestInitFromEnv(t *testing.T) {
	// Save the original environment variable and restore it after the test
	originalValue := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", originalValue)

	testCases := []struct {
		envValue string
		expected LogLevel
	}{
		{"none", LogLevelNone},
		{"error", LogLevelError},
		{"info", LogLevelInfo},
		{"debug", LogLevelDebug},
		{"invalid", LogLevelInfo}, // Default to info for invalid input
		{"", LogLevelInfo},        // Default to info if not set
	}

	for _, tc := range testCases {
		t.Run(tc.envValue, func(t *testing.T) {
			// Set the environment variable
			if tc.envValue == "" {
				os.Unsetenv("LOG_LEVEL")
			} else {
				os.Setenv("LOG_LEVEL", tc.envValue)
			}

			// Reset the log level to a known state
			SetLogLevel(LogLevelError)

			// Initialize from environment
			InitFromEnv()

			// Check the result
			result := GetLogLevel()
			if result != tc.expected {
				t.Errorf("InitFromEnv() with LOG_LEVEL=%q set log level to %v, expected %v", tc.envValue, result, tc.expected)
			}
		})
	}
}

func TestLogLevelString(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelNone, "NONE"},
		{LogLevelError, "ERROR"},
		{LogLevelInfo, "INFO"},
		{LogLevelDebug, "DEBUG"},
		{LogLevel(999), "UNKNOWN"}, // Invalid log level
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.level.String()
			if result != tc.expected {
				t.Errorf("LogLevel(%d).String() = %q, expected %q", tc.level, result, tc.expected)
			}
		})
	}
}
