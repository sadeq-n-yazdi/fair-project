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
- `pkg/docs`: API documentation (OpenAPI specification and Swagger UI)
- `pkg/errors`: Error handling and logging
- `pkg/version`: Version information

And the following entry points:

- `cmd/server`: HTTP server for the API
- `cmd/cli`: Command-line interface

## Building

To build the project, run:

```bash
cd fair-project-go
go build -o bin/fair-project-server ./cmd/server
go build -o bin/fair-project-cli ./cmd/cli
```

### Building with Git Branch Hash

To include the Git branch hash in the version information, use the following build command:

```bash
cd fair-project-go
BRANCH_HASH=$(git rev-parse --short HEAD)
go build -ldflags="-X 'github.com/sadeq/fair-project-go/pkg/version.BranchHash=$BRANCH_HASH'" -o bin/fair-project-server ./cmd/server
go build -ldflags="-X 'github.com/sadeq/fair-project-go/pkg/version.BranchHash=$BRANCH_HASH'" -o bin/fair-project-cli ./cmd/cli
```

When building with Docker, the Git branch hash is automatically included in the version information.

## Running the Server

To run the HTTP server:

```bash
./bin/fair-project-server
```

The server will start on port 8080 and provide the following endpoints:

- `GET /`: Welcome message with version information
- `GET /version`: Get the current version of the API
- `GET /docs`: API documentation (HTML UI or JSON)
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
  fair-project-cli <command> [options]

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
  help               Show this help message

Run 'fair-project-cli <command> -h' for more information on a command.
```

### Examples

Create a new class/term:

```bash
./bin/fair-project-cli create-class -name Fall2023
```

Import projects from a file:

```bash
./bin/fair-project-cli import-projects -class Fall2023 -file projects.txt
```

Import students from a file:

```bash
./bin/fair-project-cli import-students -class Fall2023 -file students.txt
```

Run the assignment algorithm:

```bash
./bin/fair-project-cli run-assignment -class Fall2023
```

Show the details of an assignment:

```bash
./bin/fair-project-cli show-assignment -class Fall2023 -id <assignment-id>
```

Export an assignment to a JSON file:

```bash
./bin/fair-project-cli export-assignment -class Fall2023 -id <assignment-id> -output assignment.json
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

Example:

```bash
docker run -p 9000:9000 -e PORT=9000 -e DATA_DIR=/data -v /path/to/data:/data fair-project-server
```

This will:
- Run the server on port 9000
- Store the data in `/data` inside the container
- Mount `/path/to/data` on your host to `/data` in the container

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
