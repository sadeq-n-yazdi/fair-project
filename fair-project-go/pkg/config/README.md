# Configuration Package

This package provides configuration settings for the Fair Project API, including log level configuration.

## Features

- Configurable log levels (none, error, info, debug)
- Thread-safe access to configuration settings
- Environment variable support for configuration

## Usage

### Log Levels

The configuration package defines the following log levels:

- `LogLevelNone`: Disables all logging
- `LogLevelError`: Logs only errors
- `LogLevelInfo`: Logs errors and informational messages (default)
- `LogLevelDebug`: Logs errors, informational messages, and debug messages

### Setting Log Level

You can set the log level in several ways:

#### 1. Using Environment Variables

Set the `LOG_LEVEL` environment variable before starting the application:

```bash
# Set log level to debug
export LOG_LEVEL=debug

# Run the application
./fair-project-server
```

Valid values for the `LOG_LEVEL` environment variable are:
- `none`
- `error`
- `info` (default)
- `debug`

#### 2. Programmatically

```go
import "github.com/sadeq/fair-project-go/pkg/config"

// Set log level to debug
config.SetLogLevel(config.LogLevelDebug)
```

### Getting the Current Log Level

```go
import "github.com/sadeq/fair-project-go/pkg/config"

// Get the current log level
level := config.GetLogLevel()
```

### Checking if a Log Level Should Be Logged

```go
import "github.com/sadeq/fair-project-go/pkg/config"

// Check if debug messages should be logged
if config.ShouldLog(config.LogLevelDebug) {
    // Log debug message
}
```

### Initializing from Environment Variables

The application automatically initializes the log level from the `LOG_LEVEL` environment variable when the `InitFromEnv` function is called:

```go
import "github.com/sadeq/fair-project-go/pkg/config"

// Initialize configuration from environment variables
config.InitFromEnv()
```

This is typically done in the `main` function of the application.
