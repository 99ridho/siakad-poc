# SIAKAD System Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Clean Architecture Implementation](#clean-architecture-implementation)
3. [Database Architecture](#database-architecture)
4. [API Design & Standards](#api-design--standards)
5. [Authentication System](#authentication-system)
6. [Module Structure](#module-structure)
7. [Infrastructure & Configuration](#infrastructure--configuration)
8. [Development Workflow](#development-workflow)
9. [Technology Stack](#technology-stack)
10. [Future Considerations](#future-considerations)

---

## System Overview

**SIAKAD (Student Information Academic System)** is an advanced proof-of-concept REST API built in Go that demonstrates production-ready patterns for managing academic systems. It includes comprehensive user authentication, complete course offering management, student enrollment system, and role-based access control.

### Current Implementation Status

- ✅ **Authentication System**: Login and registration with JWT tokens + middleware
- ✅ **Database Layer**: Complete schema with SQLC integration and soft deletes
- ✅ **API Standards**: Standardized responses with comprehensive validation
- ✅ **Clean Architecture**: Proper separation of concerns with dependency injection
- ✅ **Academic Module**: Full course enrollment system + complete course offering CRUD
- ✅ **Role-Based Access Control**: Multi-tier authorization with middleware chaining
- ✅ **Testing Framework**: Comprehensive unit tests with testify and mocking
- ✅ **Production Logging**: Structured logging with error tracking and request tracing
- ✅ **Pagination Support**: Database-level pagination with metadata

### Key Characteristics

- **Clean Architecture**: Follows Uncle Bob's clean architecture principles with clear layer separation
- **Type Safety**: SQLC-generated type-safe database queries with pgx/v5 integration
- **Comprehensive Validation**: go-playground/validator/v10 with custom error formatting
- **Advanced Security**: JWT authentication + role-based authorization with middleware chaining
- **Production Patterns**: Structured logging, error handling, and request tracking
- **Maintainability**: Modular structure with full dependency injection
- **Database Design**: Soft deletes, audit fields, and UUID primary keys
- **Test Coverage**: Comprehensive unit tests with mocking and test suites
- **Business Logic**: Advanced enrollment validation (capacity, schedule conflicts, duplicate prevention)
- **API Completeness**: Full CRUD operations with pagination and role restrictions

---

## Clean Architecture Implementation

The system implements clean architecture with four distinct layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   HTTP Handler  │  │   Validation    │                  │
│  │   (Echo v4)     │  │  (validator)    │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    USE CASE LAYER                           │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Business Logic │  │   Domain Rules  │                  │
│  │   (Use Cases)   │  │   (Validation)  │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  REPOSITORY LAYER                           │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   Interface     │  │   Implementation│                  │
│  │  (Abstraction)  │  │    (Concrete)   │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   DATABASE LAYER                            │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │      SQLC       │  │   PostgreSQL    │                  │
│  │  (Generated)    │  │     (pgx/v5)    │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### 1. Presentation Layer (`modules/*/handlers/`)

- **Purpose**: HTTP request/response handling
- **Components**: Echo handlers, request/response structs, validation
- **Dependencies**: Use cases only
- **Key Files**: `login.go`, `register.go`

#### 2. Use Case Layer (`modules/*/usecases/`)

- **Purpose**: Business logic and domain rules
- **Components**: Business logic, domain validation, orchestration
- **Dependencies**: Repository interfaces only
- **Key Files**: `login.go`, `register.go`

#### 3. Repository Layer (`db/repositories/`)

- **Purpose**: Data access abstraction
- **Components**: Interfaces and implementations
- **Dependencies**: Generated database code
- **Key Files**: `users.go`

#### 4. Database Layer (`db/generated/`)

- **Purpose**: Type-safe database operations
- **Components**: SQLC-generated code
- **Dependencies**: PostgreSQL database
- **Key Files**: `users.sql.go`, `models.go`

### Dependency Injection

All dependencies are wired in `cmd/main.go`:

```go
// Repository layer
usersRepository := repositories.NewDefaultUserRepository(pool)

// Use case layer
loginUseCase := usecases.NewLoginUseCase(usersRepository)
registerUseCase := usecases.NewRegisterUseCase(usersRepository)

// Presentation layer
loginHandler := handlers.NewLoginHandler(loginUseCase)
registerHandler := handlers.NewRegisterHandler(registerUseCase)
```

---

## Database Architecture

### Entity Relationship Diagram

```
┌─────────────┐
│    users    │
│             │
│ id (UUID)   │◄─┐
│ email       │  │
│ password    │  │
│ role        │  │
│ audit_fields│  │
└─────────────┘  │
                 │
       ┌─────────────────┐
       │                 │
┌─────────────┐  ┌───────▼──────┐
│course_reg.. │  │academic_years│
│             │  │              │
│ id (UUID)   │  │ id (UUID)    │
│ student_id  │  │ code         │
│ course_off..│  │ start_time   │
│ audit_fields│  │ end_time     │
└─────────────┘  │ audit_fields │
       ▲         └──────────────┘
       │                │
       │                ▼
┌─────────────┐  ┌─────────────┐
│course_off.. │  │ semesters   │
│             │  │             │
│ id (UUID)   │  │ id (UUID)   │
│ semester_id │◄─┤ academic_yr │
│ course_id   │  │ code        │
│ section_code│  │ start_time  │
│ capacity    │  │ end_time    │
│ start_time  │  │ audit_fields│
│ audit_fields│  └─────────────┘
└─────────────┘
       ▲
       │
┌─────────────┐
│   courses   │
│             │
│ id (UUID)   │
│ code        │
│ name        │
│ credit      │
│ audit_fields│
└─────────────┘
```

### Table Specifications

#### Users Table

```sql
CREATE TABLE users (
    id uuid not null PRIMARY KEY,
    email varchar(255) not null,
    password varchar(255) not null,
    role numeric(2) not null, -- 1=admin, 2=coordinator, 3=student
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null
);
```

**Role System**:

- `1`: Admin (full system access)
- `2`: Coordinator (program-level access)
- `3`: Student (limited access)

#### Academic Structure

- **academic_years**: Define academic periods (e.g., "2023/2024")
- **semesters**: Subdivisions within academic years (e.g., "Ganjil", "Genap")
- **courses**: Course catalog with credits
- **course_offerings**: Scheduled course sections per semester
- **course_registrations**: Student enrollment records

### SQLC Integration

#### Configuration (`sqlc.yml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./db/sql"
    schema: "./db/migrations"
    gen:
      go:
        package: "generated"
        out: "./db/generated"
        sql_package: "pgx/v5"
```

#### Query Examples (`db/sql/users.sql`)

```sql
-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: CreateUser :one
insert into users (id, email, password, role, created_at, updated_at)
values (gen_random_uuid(), $1, $2, $3, now(), now())
returning *;
```

### Migration Strategy

- **Tool**: Goose migration framework
- **Format**: Timestamped SQL files (`20250904105520_init_schema.sql`)
- **Structure**: `-- +goose Up` and `-- +goose Down` sections
- **Location**: `db/migrations/`

---

## API Design & Standards

### REST Endpoint Structure

#### Public Endpoints (Unprotected)

```
POST /login      - User authentication
POST /register   - User registration
```

#### Protected Endpoints (JWT Required)

```
# Student-only endpoints
POST /academic/course-offering/:id/enroll - Enroll student in course offering

# Admin/Coordinator-only endpoints
GET  /academic/course-offering        - List course offerings (paginated)
POST /academic/course-offering        - Create new course offering
PUT  /academic/course-offering/:id    - Update course offering
DELETE /academic/course-offering/:id  - Soft delete course offering
```

**Authentication**: Protected routes require `Authorization: Bearer <jwt-token>` header.
**Authorization**: Multi-level role-based access control:

- **Student Endpoints**: Enrollment restricted to `RoleStudent (3)` only
- **Management Endpoints**: Course offering CRUD restricted to `RoleAdmin (1)` and `RoleKoorprodi (2)`
- **Middleware Chaining**: JWT authentication + role-based authorization enforced via chained middleware

### Standardized Response Format

All API endpoints return consistent JSON responses using generic types:

```go
type BaseResponse[Data any] struct {
    Status string             `json:"status"`           // "success" | "error"
    Data   *Data              `json:"data,omitempty"`   // Response payload
    Error  *BaseResponseError `json:"error,omitempty"`  // Error details
}

type BaseResponseError struct {
    Message   string   `json:"message"`    // Human-readable error
    Details   []string `json:"details"`    // Validation errors
    Timestamp string   `json:"timestamp"`  // RFC3339 format
    Path      string   `json:"path"`       // Request URI
}
```

#### Success Response Example

```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### Error Response Example

```json
{
  "status": "error",
  "error": {
    "message": "Validation failed",
    "details": [
      "email is required",
      "password must be at least 6 characters long"
    ],
    "timestamp": "2025-01-15T10:30:00Z",
    "path": "/register"
  }
}
```

### Request Validation

Using `github.com/go-playground/validator/v10`:

```go
type LoginRequestData struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=1"`
}

type RegisterRequestData struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=6"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
```

### Validation Error Formatting

Custom error formatter in `common/validator.go` provides user-friendly messages:

- `required` → "email is required"
- `email` → "email must be a valid email address"
- `min=6` → "password must be at least 6 characters long"
- `eqfield=Password` → "confirm_password must match password"

### Pagination Support

For future list endpoints:

```go
type PaginatedBaseResponse[Data any] struct {
    BaseResponse[Data]
    Paging *PaginationMetadata `json:"paging,omitempty"`
}

type PaginationMetadata struct {
    Page         int `json:"page"`
    PageSize     int `json:"page_size"`
    TotalRecords int `json:"total_records"`
    TotalPages   int `json:"total_pages"`
}
```

---

## Authentication System

### JWT Implementation

#### Token Structure

```go
type JWTClaims struct {
    UserID string `json:"user_id"`
    Role   int64  `json:"role"`
    jwt.RegisteredClaims
}
```

#### Login Flow

1. **Validate Credentials**: Email format and required fields
2. **Authenticate User**: Verify email exists and password matches (bcrypt)
3. **Generate JWT**: 24-hour expiry with user ID and role
4. **Return Token**: Standardized success response

```go
// JWT Configuration
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
IssuedAt:  jwt.NewNumericDate(time.Now())
Subject:   user.ID.String()
```

#### Registration Flow

1. **Validate Input**: Email format, password requirements, confirmation match
2. **Check Uniqueness**: Ensure email not already registered
3. **Hash Password**: bcrypt with default cost
4. **Create User**: Default role 3 (student)
5. **Return User ID**: Success response with new user UUID

### Security Features

#### Password Security

- **Hashing**: bcrypt with default cost (currently 10)
- **Validation**: Minimum 6 characters required
- **Storage**: Only hashed passwords stored in database

#### JWT Security

- **Algorithm**: HS256 (HMAC SHA-256)
- **Secret**: Configurable via `config.json`
- **Expiry**: 24-hour token lifetime
- **Claims**: Minimal payload (user ID and role only)

#### Input Validation

- **Email**: RFC-compliant email validation
- **Password**: Length and confirmation requirements
- **Request Body**: JSON schema validation

### Role-Based Access Control

Current role hierarchy:

- **Admin (1)**: Full system administration
- **Coordinator (2)**: Program/department management
- **Student (3)**: Limited access (default for registration)

### Configuration

```json
{
  "jwt": {
    "secret": "your-secret-key-here-replace-with-secure-random-string"
  }
}
```

---

## Module Structure

### Authentication Module (`modules/auth/`)

```
modules/auth/
├── handlers/
│   ├── login.go          # POST /login endpoint
│   └── register.go       # POST /register endpoint
└── usecases/
    ├── login.go          # Login business logic
    └── register.go       # Registration business logic
```

#### Handler Layer Pattern

```go
type LoginHandler struct {
    usecase *usecases.LoginUseCase
}

func (h *LoginHandler) HandleLogin(c echo.Context) error {
    // 1. Bind request
    // 2. Validate input
    // 3. Call use case
    // 4. Return response
}
```

#### Use Case Layer Pattern

```go
type LoginUseCase struct {
    repository repositories.UserRepository
}

func (u *LoginUseCase) Login(ctx context.Context, email, password string) (string, error) {
    // 1. Business logic
    // 2. Domain validation
    // 3. Repository calls
    // 4. Return result
}
```

### Academic Module (`modules/academic/`)

**Status**: ✅ Fully implemented with course enrollment + complete course offering management

**Current Structure**:

```
modules/academic/
├── handlers/
│   ├── course_enrollment.go     # Student enrollment endpoint
│   └── course_offering.go       # Complete CRUD operations
└── usecases/
    ├── course_enrollment.go     # Enrollment business logic & validation
    ├── course_enrollment_test.go # Comprehensive enrollment unit tests
    ├── course_offering.go       # Course offering CRUD business logic
    └── course_offering_test.go  # Course offering tests
```

#### Implemented Features

**Course Enrollment System:**

- **Student Enrollment**: Role-restricted enrollment in course offerings
- **Business Validation**: Duplicate prevention, capacity checking, schedule conflict detection
- **Advanced Logic**: Time-based conflict detection with helper functions

**Course Offering Management System:**

- **Complete CRUD Operations**: Create, Read, Update, Delete (soft delete) course offerings
- **Pagination Support**: Database-level pagination with metadata
- **Role-Based Access**: Admin/Coordinator-only management operations
- **Data Integrity**: UUID handling, timestamp management, audit fields
- **Production Logging**: Comprehensive request tracking and error logging
- **Validation**: Request validation with detailed error responses

**Common Features:**

- **Error Handling**: Standardized error responses with proper HTTP status codes
- **Test Coverage**: Full unit test suites with mock dependencies
- **Production Patterns**: Structured logging, error tracking, request tracing

### Common Utilities (`common/`)

```
common/
├── base_response.go      # Standardized API responses
└── validator.go          # Request validation utilities
```

### Middleware System (`middlewares/`)

```
middlewares/
├── jwt.go               # JWT authentication middleware
└── access_control.go    # Role-based access control middleware
```

#### JWT Middleware Features (`jwt.go`)

- **Token Extraction**: Parses `Bearer <token>` from Authorization header
- **Token Validation**: Verifies JWT signature and expiration
- **Claims Extraction**: Adds user ID and role to request context
- **Error Responses**: Standardized unauthorized responses
- **Security**: HMAC SHA-256 signature verification

#### Access Control Middleware Features (`access_control.go`)

- **Role-Based Access Control**: Restricts endpoints by user roles
- **Dynamic Role Checking**: Configurable role requirements per endpoint
- **Integration**: Works seamlessly with JWT middleware
- **Authorization**: Enforces role-based authorization after authentication

### Constants (`constants/`)

```
constants/
└── constant.go          # Role definitions and system constants
```

#### Role System Definition

```go
type RoleType = int64

const (
    RoleAdmin     RoleType = 1  // System administrator
    RoleKoorprodi RoleType = 2  // Program coordinator
    RoleStudent   RoleType = 3  // Student user
)
```

---

## Infrastructure & Configuration

### Database Connection

#### PostgreSQL with pgx/v5

```go
pool, err := pgxpool.New(context.Background(), config.CurrentConfig.Database.DSN())
```

**Features**:

- Connection pooling for performance
- Context-aware operations
- Type-safe parameter binding
- Automatic connection management

#### Configuration Structure

```go
type DatabaseConfigParams struct {
    Hostname string `json:"hostname"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
    Port     string `json:"port"`
    Schema   string `json:"schema"`
}

func (c DatabaseConfigParams) DSN() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        c.Username, c.Password, c.Hostname, c.Port, c.Database,
    )
}
```

### Web Framework - Echo v4

#### Server Configuration

```go
e := echo.New()

// Public routes
e.POST("/login", loginHandler.HandleLogin)
e.POST("/register", registerHandler.HandleRegister)

// Protected routes with middleware chain
academicGroup := e.Group("/academic")
academicGroup.Use(middlewares.JWT())
academicGroup.POST(
    "/course-offering/:id/enroll",
    enrollmentHandler.HandleCourseEnrollment,
    middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
)

e.Logger.Fatal(e.Start(":8880"))
```

**Features**:

- Fast HTTP routing
- Middleware support
- Built-in JSON binding
- Context-aware request handling

### Logging - Zerolog

#### Configuration

```go
func init() {
    zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
```

**Features**:

- Structured JSON logging
- Error stack traces
- High performance
- Context-aware logging

### Configuration Management

#### File-Based Configuration (`config.json`)

```json
{
  "database": {
    "hostname": "localhost",
    "database": "siakad",
    "username": "user",
    "password": "password",
    "port": "5432",
    "schema": "public"
  },
  "jwt": {
    "secret": "your-secret-key-here"
  }
}
```

#### Configuration Loading

```go
func LoadConfig() error {
    file, err := os.ReadFile("./config.json")
    if err != nil {
        return errors.Wrap(err, "error loading config")
    }

    err = json.Unmarshal(file, &CurrentConfig)
    if err != nil {
        return errors.Wrap(err, "error loading config")
    }

    return nil
}
```

---

## Development Workflow

### Application Lifecycle

#### 1. Build and Run

```bash
# Build application
go build cmd/main.go

# Run application (development)
go run cmd/main.go

# Server starts on port 8880
```

#### 2. Configuration Setup

```bash
# Copy example configuration
cp config.json.example config.json

# Edit database credentials
vim config.json
```

### Database Operations

#### 1. Schema Migrations

```bash
# Apply migrations (requires goose)
goose -dir db/migrations postgres "connection-string" up

# Rollback migrations
goose -dir db/migrations postgres "connection-string" down
```

#### 2. Code Generation Workflow

```bash
# 1. Modify SQL queries in db/sql/*.sql
# 2. Generate type-safe Go code
sqlc generate

# Generated files appear in db/generated/
```

#### 3. Adding New Queries

1. **Define Query**: Add to `db/sql/*.sql` with SQLC annotations
2. **Generate Code**: Run `sqlc generate`
3. **Update Repository**: Add method to repository interface/implementation
4. **Update Use Case**: Use new repository method
5. **Update Handler**: Add endpoint if needed

### Dependency Management

```bash
# Add new dependency
go get github.com/package/name

# Tidy modules
go mod tidy

# Verify modules
go mod verify
```

### Testing Strategy

**Current Status**: ✅ Testify framework configured with comprehensive test coverage

**Current Implementation**:

- **Unit Tests**: Complete test suites for both enrollment and course offering use cases
- **Mock Strategy**: Repository interface mocking using testify/mock
- **Test Organization**: Structured test suites with setup/teardown methods
- **Coverage Areas**:
  - Business logic validation (enrollment rules, CRUD operations)
  - Error scenarios and edge cases
  - Helper function testing (time calculations, UUID handling)
  - Repository interaction patterns
  - Pagination logic testing

**Test Files**:

- `modules/academic/usecases/course_enrollment_test.go` - Enrollment system tests
- `modules/academic/usecases/course_offering_test.go` - Course offering CRUD tests

**POC to Production Testing Roadmap**:

- **Current POC**: Unit tests with mocks demonstrate patterns
- **Production Expansion**: Integration tests, handler tests, database test containers
- **Test Automation**: CI/CD pipeline integration for continuous testing
- **Performance Testing**: Load testing for pagination and concurrent operations

---

## Technology Stack

### Core Dependencies

#### Web Framework

- **Echo v4** (`github.com/labstack/echo/v4`): HTTP router and middleware

#### Database

- **PostgreSQL**: Primary database
- **pgx/v5** (`github.com/jackc/pgx/v5`): PostgreSQL driver with connection pooling
- **SQLC**: Type-safe SQL query generator

#### Authentication & Security

- **JWT** (`github.com/golang-jwt/jwt/v5`): Token-based authentication with middleware
- **bcrypt** (`golang.org/x/crypto/bcrypt`): Password hashing
- **Validator** (`github.com/go-playground/validator/v10`): Request validation

#### Utilities

- **Zerolog** (`github.com/rs/zerolog`): Structured logging
- **Errors** (`github.com/pkg/errors`): Enhanced error handling

#### Development Tools

- **Goose**: Database migration tool
- **SQLC**: SQL code generation

#### Testing

- **Testify** (`github.com/stretchr/testify`): Testing framework with assertions and mocks

### Directory Structure

```
siakad-poc/
├── cmd/                     # Application entry point
│   └── main.go
├── config/                  # Configuration management
│   └── config.go
├── common/                  # Shared utilities
│   ├── base_response.go     # Standardized responses
│   └── validator.go         # Request validation
├── constants/               # System constants
│   └── constant.go          # Role definitions
├── middlewares/             # HTTP middleware
│   ├── jwt.go               # JWT authentication
│   └── access_control.go    # Role-based access control
├── db/                      # Database layer
│   ├── generated/           # SQLC generated code
│   │   ├── models.go
│   │   ├── db.go
│   │   ├── users.sql.go
│   │   └── academic.sql.go
│   ├── migrations/          # Goose migration files
│   ├── repositories/        # Repository implementations
│   │   ├── users.go
│   │   └── academic.go
│   └── sql/                 # SQL query definitions
│       ├── users.sql
│       └── academic.sql
├── modules/                 # Feature modules
│   ├── auth/                # Authentication module
│   │   ├── handlers/
│   │   └── usecases/
│   └── academic/            # Academic management module
│       ├── handlers/
│       │   └── course_enrollment.go
│       └── usecases/
│           ├── course_enrollment.go
│           └── course_enrollment_test.go
├── docs/                    # Documentation
│   └── academic/
│       └── course-enrollment.md
├── config.json.example      # Configuration template
├── sqlc.yml                 # SQLC configuration
├── go.mod                   # Go module definition
├── ARCHITECTURE.md          # System architecture documentation
└── CLAUDE.md                # Development guidance
```

---

## Future Considerations

### Planned Features

#### Academic Management

- **Course Management**: CRUD operations for courses
- **Semester Management**: Academic calendar management
- **Registration System**: Student course enrollment
- **Grade Management**: Academic performance tracking

#### System Enhancements

- **Middleware**: CORS, rate limiting
- **Testing**: Comprehensive test suite
- **API Documentation**: OpenAPI/Swagger integration
- **Monitoring**: Health checks and metrics
- **CI/CD**: Automated deployment pipeline

### Scalability Considerations

#### Database

- **Read Replicas**: For query performance
- **Connection Pooling**: Already implemented
- **Indexing**: Query optimization
- **Partitioning**: For large datasets

#### Application

- **Horizontal Scaling**: Stateless design enables scaling
- **Caching**: Redis for session/data caching
- **Load Balancing**: Multiple instance deployment
- **Microservices**: Module separation for large scale

#### Security Enhancements

- **Rate Limiting**: API protection
- **CORS**: Cross-origin configuration
- **HTTPS**: TLS termination
- **Input Sanitization**: XSS protection
- **Audit Logging**: Security event tracking

### Technical Debt

#### POC Limitations (Production Refinements Needed)

**Infrastructure & Operations:**

1. **No Health Checks**: Service monitoring and readiness probes needed
2. **No Metrics Collection**: Application metrics and monitoring integration
3. **Limited Configuration**: Environment-based configuration management
4. **No Rate Limiting**: API protection and throttling mechanisms

**Security & Compliance:**

1. **Basic JWT Secret Management**: Secure secret management system needed
2. **No API Versioning**: Versioning strategy for API evolution
3. **Limited Audit Logging**: Enhanced security event tracking

**Performance & Scalability:**

1. **No Caching Layer**: Redis or similar for improved performance
2. **No Connection Pool Tuning**: Database connection optimization
3. **No Load Testing**: Performance benchmarks and capacity planning

#### Production-Ready Patterns (Already Implemented)

- ✅ **Structured Logging**: Comprehensive request tracing and error tracking
- ✅ **Role-Based Security**: Multi-tier authorization with middleware
- ✅ **Database Best Practices**: Soft deletes, audit fields, UUID keys
- ✅ **Clean Architecture**: Proper separation of concerns and dependency injection
- ✅ **Type Safety**: SQLC-generated database operations
- ✅ **Comprehensive Testing**: Unit tests with mocking patterns
- ✅ **Error Handling**: Standardized responses with detailed error information
- ✅ **Input Validation**: Request validation with custom error formatting
- ✅ **Pagination**: Database-level pagination with metadata
- ✅ **Transaction Safety**: Proper database transaction handling

---

_This document reflects the current architectural state of the SIAKAD system as of September 2025. The system represents an advanced proof-of-concept that demonstrates production-ready patterns and can be refined for production deployment. It features complete authentication, comprehensive course management (enrollment + CRUD operations), role-based access control, and extensive testing coverage. For development guidance and implementation patterns, refer to `CLAUDE.md`._
