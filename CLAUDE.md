# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a SIAKAD (Student Information Academic System) built in Go following clean architecture principles. It's a mature REST API for managing academic systems with implemented features including user authentication, JWT middleware, and a complete course enrollment system with comprehensive business validation.

## Architecture

The project follows clean architecture principles with clear separation of concerns:

- **cmd/**: Application entry point and dependency injection
- **config/**: Configuration management (JSON-based with PostgreSQL settings)  
- **common/**: Shared types and utilities (standardized API responses with generics)
- **constants/**: System constants and role definitions
- **middlewares/**: HTTP middleware for authentication and authorization
  - `jwt.go`: JWT token validation and user context injection
  - `access_control.go`: Role-based access control enforcement
- **db/**: Database layer with SQLC-generated code
  - `migrations/`: SQL migration files using goose format
  - `sql/`: SQL query definitions for SQLC
  - `generated/`: SQLC-generated Go code (models, queries)
  - `repositories/`: Repository pattern implementation
- **modules/**: Feature modules organized by domain
  - `auth/`: Authentication module with handlers and use cases
  - `academic/`: Complete course enrollment system with business logic and tests

## Key Technologies

- **Web Framework**: Echo v4 for HTTP routing and middleware
- **Database**: PostgreSQL with pgx/v5 driver and connection pooling
- **Code Generation**: SQLC for type-safe database queries
- **Migrations**: Goose (based on migration file naming)
- **Logging**: Zerolog with structured logging and error stack traces
- **Testing**: Testify framework with mocks and test suites
- **Authentication**: Separate JWT authentication and role-based access control middleware

## Database Schema

The system models academic entities with proper relationships:
- `users` (with role-based access: 1=admin, 2=coordinator, 3=student)  
- `academic_years` and `semesters` (hierarchical time periods)
- `courses` and `course_offerings` (course catalog and scheduled sections)
- `course_registrations` (student enrollment tracking)

All tables use UUID primary keys and include audit fields (created_at, updated_at, deleted_at).

## Development Commands

### Running the Application
```bash
go run cmd/main.go
```
The server starts on port 8880 by default.

### Role-Based Access Control
The system implements a three-tier role hierarchy defined in `constants/constant.go`:
- `RoleAdmin (1)`: System administrator with full access
- `RoleKoorprodi (2)`: Program coordinator with management access
- `RoleStudent (3)`: Student with limited access to enrollment features

Endpoints are protected using middleware chains:
```go
// Example: Student-only enrollment endpoint
academicGroup.POST(
    "/course-offering/:id/enroll",
    enrollmentHandler.HandleCourseEnrollment,
    middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
)
```

### Configuration
Copy `config.json.example` to `config.json` and update database credentials.

### Database Operations
```bash
# Generate SQLC code after modifying db/sql/*.sql files
sqlc generate

# Run tests to ensure changes don't break existing functionality
go test ./...

# Run migrations (if goose is installed)
goose -dir db/migrations postgres "your-connection-string" up
```

### Code Generation
After adding new SQL queries to `db/sql/`, run `sqlc generate` to regenerate the type-safe Go code, then run tests to verify integration.

### Testing
Run `go test ./...` to execute all tests. New features should include comprehensive unit tests following the existing patterns in `modules/academic/usecases/course_enrollment_test.go`.

## Architecture Patterns

### Repository Pattern
Database access is abstracted through repository interfaces, making testing and mocking easier.

### Use Case Pattern  
Business logic is encapsulated in use case structs that depend on repository interfaces.

### Handler Pattern
HTTP handlers are thin layers that handle request/response marshaling and call use cases.

### Middleware Pattern
Centralized cross-cutting concerns through Echo middleware:
- **JWT Authentication**: Token validation and user context injection (`middlewares/jwt.go`)
- **Access Control**: Role-based authorization enforcement (`middlewares/access_control.go`)
- **Chained Middleware**: Combined authentication and authorization for protected routes

### Dependency Injection
Dependencies are wired up in `cmd/main.go` following constructor injection pattern.

### Error Handling
All handlers return standardized JSON responses using `common.BaseResponse` with proper HTTP status codes and error details.

## API Standards

- REST endpoints with proper HTTP methods
- Standardized JSON responses with `status`, `data`, and `error` fields  
- Generic response types for type safety
- Pagination support via `PaginatedBaseResponse`
- Consistent error response format with timestamps and request paths
- JWT-protected routes requiring `Authorization: Bearer <token>` header
- Public authentication endpoints (`/login`, `/register`)
- Protected academic endpoints (`/academic/*` routes) with role-based access control
- Role-restricted enrollment endpoint (students only)

## Testing

### Current Testing Setup
- **Framework**: Testify with test suites and mocks
- **Test Location**: `modules/academic/usecases/course_enrollment_test.go`
- **Coverage**: Comprehensive unit tests for business logic
- **Mocking**: Repository layer mocked using testify/mock
- **Test Patterns**: Test suites with setup/teardown methods

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test suite
go test -v ./modules/academic/usecases/
```

### Test Structure
- **Business Logic Tests**: Enrollment validation, capacity checks, schedule conflicts
- **Error Scenario Tests**: Repository errors, invalid data, business rule violations
- **Helper Function Tests**: Time calculations, overlap detection
- **Mock Setup**: Repository mocks with expected calls and returns