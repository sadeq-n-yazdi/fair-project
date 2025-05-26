# API Testing Documentation

This document describes the testing approach and results for the Fair Project Go API.

## Testing Approach

The API endpoints are tested using Go's built-in testing framework. The tests are organized as follows:

1. **Test Utilities**: Common functions for setting up test environments, creating test requests, and checking responses.
2. **HTTP Utilities Tests**: Tests for the HTTP utility functions that handle JSON responses.
3. **Authentication Tests**: Tests for the authentication-related endpoints.
4. **Class Handlers Tests**: Tests for the class-related endpoints.
5. **Assignment Handlers Tests**: Tests for the assignment-related endpoints.
6. **Server Tests**: Tests for the server's routing logic.

### Test Environment

The tests use a temporary directory for storing data, which is created at the start of each test and cleaned up afterward. This ensures that tests don't interfere with each other or with the actual data directory.

### Test Utilities

The test utilities include functions for:

- Setting up a test environment with a temporary data directory
- Creating HTTP requests with JSON bodies
- Executing requests and capturing responses
- Checking response codes and content types
- Parsing response bodies
- Creating test data (classes, projects, students)

## Tested Endpoints

The following API endpoints have been tested:

1. **GET /**: Returns a welcome message.
2. **POST /auth/login**: Authenticates a user and returns a JWT token.
3. **GET /auth/users**: Lists all users (admin and superadmin only).
4. **POST /auth/users**: Creates a new user (admin and superadmin only).
5. **GET /auth/users/{username}**: Gets a specific user (admin and superadmin only, or the user themselves).
6. **PUT /auth/users/{username}**: Updates a user's roles and enabled status (admin and superadmin only).
7. **DELETE /auth/users/{username}**: Deletes a user (superadmin only).
8. **GET /user/whoami**: Gets the current user's information (username, status, and roles).
9. **POST /classes**: Creates a new class/term.
10. **GET /classes**: Lists all class/terms.
11. **POST /classes/{className}/projects**: Uploads projects for a class/term.
12. **GET /classes/{className}/projects**: Gets projects for a class/term.
13. **POST /classes/{className}/students**: Uploads students for a class/term.
14. **GET /classes/{className}/students**: Gets students for a class/term.
15. **POST /classes/{className}/assign**: Runs the assignment algorithm for a class/term.
16. **GET /classes/{className}/assignments**: Lists all assignments for a class/term.
17. **GET /classes/{className}/assignments/{assignmentID}**: Gets a specific assignment result.

## Test Results

All tests have passed successfully, confirming that the API endpoints are working correctly.

```
ok      github.com/sadeq/fair-project-go/pkg/api        0.004s
ok      github.com/sadeq/fair-project-go/cmd/server     0.002s
```

## Test Coverage

The tests cover the following aspects of the API:

- **Input Validation**: Tests include valid and invalid inputs to ensure proper validation.
- **Error Handling**: Tests verify that appropriate error responses are returned for invalid inputs or when resources don't exist.
- **Success Cases**: Tests verify that successful operations return the expected responses.
- **End-to-End Workflows**: Tests simulate complete workflows, such as creating a class, uploading projects and students, running an assignment, and retrieving the results.

## Future Improvements

Potential improvements to the testing approach include:

1. **Mocking the Storage Layer**: Currently, the tests use a real file system with a temporary directory. Mocking the storage layer would make the tests faster and more isolated.
2. **Property-Based Testing**: Using property-based testing to generate random inputs and verify that the API behaves correctly for a wide range of inputs.
3. **Load Testing**: Testing the API under load to ensure it can handle multiple concurrent requests.
4. **Integration Tests**: Testing the API in conjunction with a real client to ensure end-to-end functionality.
