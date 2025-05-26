# Fair Project Go

A tool for assigning projects to students based on their preferences in a fair way.

**Current Version:** 0.0.1 (includes Git branch hash when built with Docker or using build flags)

## Project Structure

The project has been restructured into the following packages:

- `pkg/models`: Data structures used throughout the application
- `pkg/assignment`: Core assignment logic
- `pkg/storage`: Data persistence and file I/O operations
- `pkg/api`: HTTP API handlers
- `pkg/cli`: Command-line interface functions
- `pkg/config`: Configuration settings for the application
- `pkg/docs`: API documentation (OpenAPI specification and Swagger UI)
- `pkg/errors`: Error handling and logging
- `pkg/middleware`: HTTP middleware components
  - `logging`: Request/response logging middleware
- `pkg/version`: Version information

And the following entry points:

- `cmd/server`: HTTP server for the API
- `cmd/cli`: Command-line interface

## Building

There are several ways to build the project:

### Using Make

The project includes a Makefile with various targets for building, testing, and running the application:

```bash
cd fair-project-go

# Build both server and CLI binaries
make

# Build only the server binary
make build-server

# Build only the CLI binary
make build-cli

# Run tests
make test

# Clean build artifacts
make clean

# Build Docker image
make docker

# Show all available targets
make help
```

### Using the Build Script

Alternatively, you can use the build script:

```bash
cd fair-project-go

# Make the script executable
chmod +x scripts/build.sh

# Build both server and CLI binaries
./scripts/build.sh

# Build only the server binary
./scripts/build.sh server

# Build only the CLI binary
./scripts/build.sh cli

# Run tests
./scripts/build.sh test

# Clean build artifacts
./scripts/build.sh clean

# Show all available commands
./scripts/build.sh help
```

### Manual Building

You can also build the project manually:

```bash
cd fair-project-go

# Build without commit hash
go build -o bin/fair-project-server ./cmd/server
go build -o bin/fair-ctl ./cmd/cli

# Build with commit hash
COMMIT_HASH=$(git rev-parse --short HEAD)
go build -ldflags="-X 'github.com/sadeq/fair-project-go/pkg/version.BranchHash=$COMMIT_HASH'" -o bin/fair-project-server ./cmd/server
go build -ldflags="-X 'github.com/sadeq/fair-project-go/pkg/version.BranchHash=$COMMIT_HASH'" -o bin/fair-ctl ./cmd/cli
```

### Docker Building

When building with Docker, the Git commit hash is automatically included in the version information:

```bash
cd fair-project-go
docker build -t fair-project-server .
```

## Running the Server

To run the HTTP server:

```bash
./bin/fair-project-server
```

The server will start on port 8080 and provide the following endpoints:

- `GET /`: Welcome message with version information
- `GET /version`: Get the current version of the API
- `GET /docs`: API documentation (HTML UI or JSON)
- `POST /auth/login`: Authenticate a user and get a JWT token
- `GET /auth/users`: List all users (admin and superadmin only)
- `POST /auth/users`: Create a new user (admin and superadmin only)
- `GET /auth/users/{username}`: Get a specific user (admin and superadmin only, or the user themselves)
- `PUT /auth/users/{username}`: Update a user's roles and enabled status (admin and superadmin only)
- `DELETE /auth/users/{username}`: Delete a user (superadmin only)
- `GET /user/whoami`: Get the current user's information (username, status, and roles)
- `POST /classes`: Create a new class/term
- `GET /classes`: List all class/terms
- `POST /classes/{className}/projects`: Upload projects for a class/term
- `GET /classes/{className}/projects`: Get projects for a class/term
- `POST /classes/{className}/students`: Upload students for a class/term
- `GET /classes/{className}/students`: Get students for a class/term
- `POST /classes/{className}/assign`: Run the assignment algorithm for a class/term
- `GET /classes/{className}/assignments`: List all assignments for a class/term
- `GET /classes/{className}/assignments/{assignmentID}`: Get a specific assignment result

## Using the CLI

The CLI provides a command-line interface for the application. Here are the available commands:

```
Usage:
  fair-ctl <command> [options]

Commands:
  create-class       Create a new class/term
  list-classes       List all class/terms
  import-projects    Import projects from a file
  import-students    Import students from a file
  list-projects      List all projects for a class/term
  list-students      List all students for a class/term
  run-assignment     Run the assignment algorithm for a class/term
  list-assignments   List all assignments for a class/term
  show-assignment    Show the details of an assignment
  export-assignment  Export an assignment to a JSON file
  create-superadmin  Create the first superadmin user
  change-password    Change a user's password
  list-users         List all users
  generate-env       Generate a new .env file with random values for sensitive fields
  completion         Generate shell completion scripts
  help               Show this help message

Run 'fair-ctl <command> -h' for more information on a command.
```

### Examples

Create a new class/term:

```bash
./bin/fair-ctl create-class -name Fall2023
```

Import projects from a file:

```bash
./bin/fair-ctl import-projects -class Fall2023 -file projects.txt
```

Import students from a file:

```bash
./bin/fair-ctl import-students -class Fall2023 -file students.txt
```

Run the assignment algorithm:

```bash
./bin/fair-ctl run-assignment -class Fall2023
```

Show the details of an assignment:

```bash
./bin/fair-ctl show-assignment -class Fall2023 -id <assignment-id>
```

Export an assignment to a JSON file:

```bash
./bin/fair-ctl export-assignment -class Fall2023 -id <assignment-id> -output assignment.json
```

Generate a new .env file with random values for sensitive fields:

```bash
./bin/fair-ctl generate-env -output .env
```

This will create a new .env file with random values for JWT_SECRET, SUPERADMIN_KEY, and POSTGRES_PASSWORD, while setting
sensible defaults for other configuration values.

### Shell Autocompletion

The CLI supports shell autocompletion for bash, zsh, and fish shells. This feature helps you quickly complete commands,
flags, and arguments by pressing the Tab key.

#### Generating Completion Scripts

You can generate shell completion scripts using the `completion` command:

```bash
# For Bash
./bin/fair-ctl completion bash > ~/.fair-ctl-completion.bash

# For Zsh
./bin/fair-ctl completion zsh > ~/.fair-ctl-completion.zsh

# For Fish
./bin/fair-ctl completion fish > ~/.config/fish/fair-ctl-completion.fish
```

You can also specify an output file directly:

```bash
./bin/fair-ctl completion -shell bash -output ~/.fair-ctl-completion.bash
```

#### Activating Autocompletion

After generating the completion script, you need to source it in your shell configuration file:

**Bash**:

```bash
echo "source ~/.fair-ctl-completion.bash" >> ~/.bashrc
source ~/.bashrc
```

**Zsh**:

```bash
echo "source ~/.fair-ctl-completion.zsh" >> ~/.zshrc
source ~/.zshrc
```

**Fish**:

```bash
# Fish automatically loads completions from ~/.config/fish/
source ~/.config/fish/fair-ctl-completion.fish
```

#### One-Step Installation

You can also use the CLI to automatically install and activate completion for your current shell:

```bash
./bin/fair-ctl completion -shell $(basename $SHELL) -output ~/.$(basename $SHELL)rc_fair-ctl
```

This will generate the appropriate completion script and add it to your shell configuration file.

## Logging

The API server includes a logging middleware that logs all HTTP requests and responses. The verbosity of the logging can be configured using the `LOG_LEVEL` environment variable:

- `none`: Disables all logging
- `error`: Logs only errors
- `info`: Logs errors and informational messages (default)
- `debug`: Logs errors, informational messages, and debug messages (including request/response bodies)

Example log output at the `info` level:

```
2023/05/01 12:34:56 [INFO] Request: GET /classes
2023/05/01 12:34:56 [INFO] Response: GET /classes - Status: 200 - Duration: 5.123ms
```

Example log output at the `debug` level:

```
2023/05/01 12:34:56 [INFO] Request: POST /classes
2023/05/01 12:34:56 [DEBUG] Request Headers: map[Content-Type:[application/json] User-Agent:[curl/7.68.0]]
2023/05/01 12:34:56 [DEBUG] Request Body: {"name":"Fall2023"}
2023/05/01 12:34:56 [INFO] Response: POST /classes - Status: 201 - Duration: 10.456ms
2023/05/01 12:34:56 [DEBUG] Response Body: {"id":"123","name":"Fall2023"}
```

## API Documentation

The API is documented using the OpenAPI 3.0 specification. You can access the documentation in two ways:

1. **Web UI**: Visit the `/docs` endpoint in your web browser to view the interactive API documentation using Swagger UI.
2. **JSON Format**: To get the OpenAPI specification in JSON format, you can:
   - Add the query parameter `?format=json` to the `/docs` endpoint
   - Set the `Accept: application/json` header in your request to the `/docs` endpoint

The documentation includes all API endpoints, request/response formats, and data models.

## Docker

The API server can be run in a Docker container. The Dockerfile is provided in the repository.

### Building the Docker Image

To build the Docker image, run:

```bash
cd fair-project-go
docker build -t fair-project-server .
```

Or use the Makefile:

```bash
make docker
```

### Running the Docker Container

To run the Docker container:

```bash
docker run -p 8080:8080 -v /path/to/data:/app/data fair-project-server
```

This will:
- Map port 8080 on your host to port 8080 in the container
- Mount `/path/to/data` on your host to `/app/data` in the container (where the data will be stored)

### Configuration

The Docker container can be configured using environment variables:

- `PORT`: The port on which the server will listen (default: 8080)
- `DATA_DIR`: The directory where the data will be stored (default: /app/data)
- `LOG_LEVEL`: The verbosity level of logging (default: info)
  - `none`: Disables all logging
  - `error`: Logs only errors
  - `info`: Logs errors and informational messages
  - `debug`: Logs errors, informational messages, and debug messages (including request/response bodies)
- `JWT_SECRET`: The secret key used to sign JWT tokens (default: a random string)
- `SUPERADMIN_KEY`: The key used for creating the first super admin via CLI

#### Using Environment Variables Directly

You can pass environment variables directly to the Docker container:

```bash
docker run -p 9000:9000 -e PORT=9000 -e DATA_DIR=/data -e LOG_LEVEL=debug -v /path/to/data:/data fair-project-server
```

This will:
- Run the server on port 9000
- Store the data in `/data` inside the container
- Set the log level to debug (verbose logging)
- Mount `/path/to/data` on your host to `/data` in the container

#### Using an Environment File

You can also use an environment file (`.env`) to configure the Docker container:

1. Create a `.env` file with your configuration (or copy and modify the provided `.env.example`):

```
# Server port
PORT=9000

# Data directory
DATA_DIR=/data

# Log level
LOG_LEVEL=debug

# JWT Secret
JWT_SECRET=your_jwt_secret_here

# Other configuration...
```

2. Run the Docker container with the `--env-file` option:

```bash
docker run -p 9000:9000 --env-file .env -v /path/to/data:/data fair-project-server
```

## Docker Compose

For a more complete setup, you can use Docker Compose to run the API server along with a PostgreSQL database.

### Starting the Services

To start all services:

```bash
cd fair-project-go
docker-compose up -d
```

Or use the Makefile:

```bash
make docker-compose-up
```

This will start:

- The API server on port 8080
- A PostgreSQL database on port 5432

### Stopping the Services

To stop all services:

```bash
cd fair-project-go
docker-compose down
```

Or use the Makefile:

```bash
make docker-compose-down
```

### Viewing Logs

To view the logs from all services:

```bash
cd fair-project-go
docker-compose logs -f
```

Or use the Makefile:

```bash
make docker-compose-logs
```

### Configuration

The services can be configured using environment variables. Docker Compose is set up to automatically use variables from
a `.env` file in the project root.

#### Using the .env File

The project includes an `.env.example` file that you can copy to create your own `.env` file:

```bash
cp .env.example .env
```

Then edit the `.env` file to set your configuration:

```
# JWT Secret for authentication
JWT_SECRET=your_jwt_secret_here

# Super Admin Key
SUPERADMIN_KEY=your_superadmin_key_here

# Server port
PORT=8080

# Data directory
DATA_DIR=data

# Log level
LOG_LEVEL=INFO

# PostgreSQL Configuration
POSTGRES_USER=fairuser
POSTGRES_PASSWORD=your_postgres_password_here
POSTGRES_DB=fairdb
POSTGRES_PORT=5432
```

Docker Compose will automatically use these variables when you run `docker-compose up`.

#### API Server Configuration

The API server can be configured using the following environment variables:

- `PORT`: The port on which the server will listen (default: 8080)
- `DATA_DIR`: The directory where the data will be stored (default: /app/data)
- `LOG_LEVEL`: The verbosity level of logging (default: info)
- `JWT_SECRET`: The secret key used to sign JWT tokens
- `SUPERADMIN_KEY`: The key used for creating the first super admin via CLI

#### PostgreSQL Configuration

The PostgreSQL database can be configured using the following environment variables:

- `POSTGRES_USER`: The username for the database (default: fairuser)
- `POSTGRES_PASSWORD`: The password for the database (default: fairpassword)
- `POSTGRES_DB`: The name of the database (default: fairdb)
- `POSTGRES_PORT`: The port on which PostgreSQL will listen (default: 5432)

#### Overriding Environment Variables

You can also override environment variables when starting Docker Compose:

```bash
PORT=9000 POSTGRES_PORT=5433 docker-compose up -d
```

This will use the values from the command line instead of the `.env` file for those specific variables.

## File Formats

### Projects File

The projects file should be a text file with one project name per line. For example:

```
Project1
Project2
Project3
```

### Students File

The students file should be a CSV file with the following format:

```
StudentName1,Preference1,Preference2,Preference3
StudentName2,Preference1,Preference2,Preference3
```

Each line represents a student, with the first column being the student's name and the remaining columns being their project preferences in order of preference.
