package config

import (
	"os"
	"strings"
	"sync"
)

// LogLevel represents the verbosity level of logging
type LogLevel int

const (
	// LogLevelNone disables all logging
	LogLevelNone LogLevel = iota
	// LogLevelError logs only errors
	LogLevelError
	// LogLevelInfo logs errors and informational messages
	LogLevelInfo
	// LogLevelDebug logs errors, informational messages, and debug messages
	LogLevelDebug
)

var (
	currentLogLevel LogLevel = LogLevelInfo // Default log level
	logLevelMutex   sync.RWMutex
)

// SetLogLevel sets the current log level
func SetLogLevel(level LogLevel) {
	logLevelMutex.Lock()
	defer logLevelMutex.Unlock()
	currentLogLevel = level
}

// GetLogLevel returns the current log level
func GetLogLevel() LogLevel {
	logLevelMutex.RLock()
	defer logLevelMutex.RUnlock()
	return currentLogLevel
}

// ShouldLog returns true if the given log level should be logged
// based on the current log level setting
func ShouldLog(level LogLevel) bool {
	logLevelMutex.RLock()
	defer logLevelMutex.RUnlock()
	return level <= currentLogLevel
}

// LogLevelFromString converts a string to a LogLevel
func LogLevelFromString(level string) LogLevel {
	switch strings.ToLower(level) {
	case "none":
		return LogLevelNone
	case "error":
		return LogLevelError
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	default:
		return LogLevelInfo // Default to info level
	}
}

// InitFromEnv initializes the log level from the LOG_LEVEL environment variable
func InitFromEnv() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		SetLogLevel(LogLevelFromString(logLevel))
	} else {
		SetLogLevel(LogLevelInfo) // Set to default log level if not specified
	}
}

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelNone:
		return "NONE"
	case LogLevelError:
		return "ERROR"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
