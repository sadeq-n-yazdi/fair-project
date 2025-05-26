# Fair Project Go

A tool for assigning projects to students based on their preferences in a fair way.

## Project Structure

The project has been restructured into the following packages:

- `pkg/models`: Data structures used throughout the application
- `pkg/assignment`: Core assignment logic
- `pkg/storage`: Data persistence and file I/O operations
- `pkg/api`: HTTP API handlers
- `pkg/cli`: Command-line interface functions

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

## Running the Server

To run the HTTP server:

```bash
./bin/fair-project-server
```

The server will start on port 8080 and provide the following endpoints:

- `GET /`: Welcome message
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
