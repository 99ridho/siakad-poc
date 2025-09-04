# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a SIAKAD (Student Information Academic System) built in Go following clean architecture principles. It's an advanced proof-of-concept REST API demonstrating production-ready patterns for managing academic systems. Features include user authentication, JWT middleware, complete course enrollment system, full course offering CRUD operations, role-based access control, and comprehensive business validation.

## Architecture

The project follows clean architecture principles with clear separation of concerns:

- **cmd/**: Application entry point and dependency injection
- **config/**: Configuration management (JSON-based with PostgreSQL settings)
- **common/**: Shared types and utilities (standardized API responses with generics)
- **constants/**: System constants and role definitions
- **middlewares/**: HTTP middleware for authentication and authorization
  - `jwt.go`: JWT token validation and user context injection with comprehensive logging
  - `access_control.go`: Role-based access control enforcement with middleware chaining
- **db/**: Database layer with SQLC-generated code
  - `migrations/`: SQL migration files using goose format
  - `sql/`: SQL query definitions for SQLC
  - `generated/`: SQLC-generated Go code (models, queries)
  - `repositories/`: Repository pattern implementation
- **modules/**: Feature modules organized by domain
  - `auth/`: Authentication module with handlers and use cases
  - `academic/`: Complete academic management system:
    - Course enrollment with advanced business validation
    - Course offering CRUD operations with pagination
    - Role-based endpoint access control
    - Comprehensive test coverage

## Key Technologies

- **Web Framework**: Fiber v2 for HTTP routing, middleware, and request handling
- **Database**: PostgreSQL with pgx/v5 driver, connection pooling, and UUID primary keys
- **Code Generation**: SQLC for type-safe database queries and model generation
- **Migrations**: Goose for database schema versioning (timestamp-based files)
- **Logging**: Zerolog with structured JSON logging, error stack traces, and request tracking
- **Testing**: Testify framework with assertion helpers, mocks, and organized test suites
- **Security**: JWT authentication + role-based access control with Fiber middleware chaining
- **Validation**: go-playground/validator/v10 with custom error formatting

## Database Schema

The system models academic entities with proper relationships:

- `users` (with role-based access: 1=admin, 2=coordinator, 3=student)
- `academic_years` and `semesters` (hierarchical time periods)
- `courses` and `course_offerings` (course catalog and scheduled sections)
- `course_registrations` (student enrollment tracking)

All tables use UUID primary keys and include comprehensive audit fields (created_at, updated_at, deleted_at) with soft delete functionality.

## Development Commands

### Running the Application

```bash
go run cmd/main.go
```

The server starts on port 8880 by default.

### Role-Based Access Control

The system implements a comprehensive three-tier role hierarchy defined in `constants/constant.go`:

- `RoleAdmin (1)`: System administrator with full system access
- `RoleKoorprodi (2)`: Program coordinator with course management access
- `RoleStudent (3)`: Student with limited access to enrollment features

Endpoints are protected using chained Fiber middleware for authentication + authorization:

```go
// Student-only enrollment endpoint
academicGroup.Post(
    "/course-offering/:id/enroll",
    enrollmentHandler.HandleCourseEnrollment,
    middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
)

// Admin/Coordinator course management endpoints
academicGroup.Get(
    "/course-offering",
    courseOfferingHandler.HandleListCourseOfferings,
    middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
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

# Run tests with verbose output for detailed results
go test -v ./...

# Run specific test suite (example: academic module)
go test -v ./modules/academic/usecases/

# Run migrations (if goose is installed)
goose -dir db/migrations postgres "your-connection-string" up
```

### Code Generation

After adding new SQL queries to `db/sql/`, run `sqlc generate` to regenerate the type-safe Go code. Follow the existing patterns for:

- UUID handling in repository methods
- Soft delete implementation (use `deleted_at IS NULL` in WHERE clauses)
- Pagination with LIMIT/OFFSET
- Comprehensive error handling with proper error messages
  Then run tests to verify integration and add corresponding unit tests.

### Testing

Run `go test ./...` to execute all tests. New features should include comprehensive unit tests following the existing patterns:

- Course enrollment tests: `modules/academic/usecases/course_enrollment_test.go`
- Course offering CRUD tests: `modules/academic/usecases/course_offering_test.go`

Test patterns to follow:

- Use testify/suite for organized test structure
- Mock repository interfaces using testify/mock
- Include business logic validation tests
- Test error scenarios and edge cases
- Verify proper error propagation and handling

## Architecture Patterns

### Repository Pattern

Database access is abstracted through repository interfaces, making testing and mocking easier.

### Use Case Pattern

Business logic is encapsulated in use case structs that depend on repository interfaces.

### Handler Pattern

HTTP handlers are thin layers that handle request/response marshaling and call use cases.

### Middleware Pattern

Centralized cross-cutting concerns through Fiber middleware:

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
- **Public Endpoints**: Authentication routes (`/login`, `/register`)
- **Protected Academic Endpoints**: Complete `/academic/*` route group with JWT authentication
- **Role-Based Access Control**:
  - Student-only: Course enrollment (`POST /academic/course-offering/:id/enroll`)
  - Admin/Coordinator-only: Course offering CRUD operations
    - `GET /academic/course-offering` (list with pagination)
    - `POST /academic/course-offering` (create)
    - `PUT /academic/course-offering/:id` (update)
    - `DELETE /academic/course-offering/:id` (soft delete)

## Testing

### Current Testing Setup

- **Framework**: Testify with assertion helpers, test suites, and comprehensive mocking
- **Test Locations**:
  - `modules/academic/usecases/course_enrollment_test.go` - Enrollment system tests
  - `modules/academic/usecases/course_offering_test.go` - Course offering CRUD tests
- **Coverage**:
  - Business logic validation (enrollment rules, CRUD operations)
  - Error scenarios and edge cases
  - Repository interaction patterns
  - Helper function testing (time calculations, UUID conversion)
- **Mocking Strategy**: Repository interface mocks using testify/mock
- **Test Organization**: Structured test suites with setup/teardown methods and grouped test cases

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

**Course Enrollment Tests:**

- Business logic validation (duplicate prevention, capacity limits)
- Schedule conflict detection with time overlap calculations
- Error scenarios (repository failures, invalid data)
- Helper function testing (time calculations, conflict detection)

**Course Offering CRUD Tests:**

- CRUD operation validation (create, read, update, delete)
- Pagination logic testing
- UUID handling and conversion testing
- Repository interaction patterns

**Mock Strategy:**

- Repository interface mocks with expected method calls and returns
- Context-aware testing with proper error propagation
- Test data factories for consistent test setup

## Production Readiness Assessment

### Current POC Status

This system demonstrates **advanced POC patterns** that can be refined for production:

**Production-Ready Patterns (✅ Implemented):**

- Clean architecture with proper dependency injection
- Comprehensive JWT authentication + role-based authorization
- Type-safe database operations with SQLC
- Structured logging with request tracking and error stack traces
- Soft delete functionality with audit fields
- Pagination support with database-level optimization
- Comprehensive input validation with detailed error responses
- Unit testing patterns with mocking strategies

**Production Refinements Needed:**

- Health check endpoints for monitoring and orchestration
- Environment-based configuration management
- Metrics collection and monitoring integration
- Rate limiting and API protection mechanisms
- Enhanced security (secret management, audit logging)
- Caching layer (Redis) for improved performance
- Load testing and capacity planning
- CI/CD pipeline integration

### Web Framework Migration: Echo v4 → Fiber v2

The system has been successfully migrated from Echo v4 to Fiber v2, maintaining all existing functionality while improving performance and developer experience:

**Migration Benefits:**
- **Performance**: Fiber v2 offers significantly better performance than Echo v4
- **Express.js-like API**: More familiar patterns for developers from Node.js background
- **Built-in Features**: Rich set of built-in middleware and utilities
- **Active Development**: Regular updates and strong community support
- **Memory Efficiency**: Lower memory footprint and faster request processing

**Key Changes:**
- **Server Setup**: `echo.New()` → `fiber.New()`
- **Route Methods**: `app.POST()` → `app.Post()`, `app.GET()` → `app.Get()`
- **Handler Signatures**: `func(c echo.Context) error` → `func(c *fiber.Ctx) error`
- **Request Handling**: `c.Bind()` → `c.BodyParser()`, `c.Param()` → `c.Params()`
- **Response Methods**: `c.JSON()` → `c.Status().JSON()`, `c.RealIP()` → `c.IP()`

**What Remained Unchanged:**
- ✅ Clean architecture and business logic
- ✅ Database operations and SQLC integration
- ✅ JWT authentication and authorization
- ✅ Request validation and error handling
- ✅ Testing patterns and frameworks
- ✅ All API endpoints and contracts

### Development Guidance for Production

**Extending Academic Features:**

1. Follow existing patterns in `modules/academic/` for new features
2. Implement comprehensive logging using the established patterns
3. Add role-based access control using middleware chaining
4. Include pagination for list endpoints
5. Write comprehensive unit tests with mocks
6. Follow the UUID handling patterns in repositories

**API Development Standards:**

1. Use standardized response format from `common/base_response.go`
2. Implement proper validation using validator tags
3. Include structured logging with request tracking
4. Handle errors with appropriate HTTP status codes
5. Document role requirements for protected endpoints

# important-instruction-reminders

- Do what has been asked; nothing more, nothing less.
- NEVER create files unless they're absolutely necessary for achieving your goal.
- ALWAYS prefer editing an existing file to creating a new one.
- NEVER proactively create documentation files (\*.md) or README files. Only create documentation files if explicitly requested by the User.
- Track your works using TODO list, so you didn't get lost.
