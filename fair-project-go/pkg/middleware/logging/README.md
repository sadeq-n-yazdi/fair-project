# Logging Middleware

This package provides middleware for logging HTTP requests and responses in the Fair Project API.

## Features

- Logs HTTP request details (method, path, headers, body)
- Logs HTTP response details (status code, body, duration)
- Configurable log levels (none, error, info, debug)
- Minimal performance impact when logging is disabled

## Usage

### In HTTP Server

```go
import (
    "net/http"
    "github.com/sadeq/fair-project-go/pkg/middleware/logging"
)

// Create a handler with the logging middleware
handler := logging.Middleware(http.HandlerFunc(yourHandler))

// Use the handler with the HTTP server
http.ListenAndServe(":8080", handler)
```

### With http.HandleFunc

```go
import (
    "net/http"
    "github.com/sadeq/fair-project-go/pkg/middleware/logging"
)

// Create a handler function with the logging middleware
handlerFunc := logging.MiddlewareFunc(yourHandlerFunc)

// Register the handler function
http.HandleFunc("/path", handlerFunc)
```

## Configuration

The logging middleware uses the configuration package to determine the log level. The log level can be set using the `LOG_LEVEL` environment variable or programmatically using the `config.SetLogLevel` function.

### Log Levels

- `none`: Disables all logging
- `error`: Logs only errors
- `info`: Logs errors and informational messages (default)
- `debug`: Logs errors, informational messages, and debug messages (including request/response bodies)

### Setting Log Level via Environment Variable

```bash
# Set log level to debug
export LOG_LEVEL=debug

# Run the server
./fair-project-server
```

### Setting Log Level Programmatically

```go
import "github.com/sadeq/fair-project-go/pkg/config"

// Set log level to debug
config.SetLogLevel(config.LogLevelDebug)
```
