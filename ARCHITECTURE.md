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

**SIAKAD (Student Information Academic System)** is a proof-of-concept REST API built in Go that manages academic systems with features including user authentication, course registration, and semester management.

### Current Implementation Status
- ✅ **Authentication System**: Login and registration with JWT tokens + middleware
- ✅ **Database Layer**: Complete schema with SQLC integration
- ✅ **API Standards**: Standardized responses with validation
- ✅ **Clean Architecture**: Proper separation of concerns
- ✅ **Academic Module**: Full course enrollment system implemented
- ✅ **Testing Framework**: Comprehensive unit tests with testify
- ✅ **JWT Middleware**: Centralized authentication for protected routes

### Key Characteristics
- **Clean Architecture**: Follows Uncle Bob's clean architecture principles
- **Type Safety**: SQLC-generated type-safe database queries
- **Validation**: go-playground/validator/v10 for request validation
- **Security**: JWT authentication with bcrypt password hashing + middleware
- **Maintainability**: Modular structure with dependency injection
- **Test Coverage**: Comprehensive unit tests with mocking
- **Business Logic**: Advanced enrollment validation (capacity, schedule conflicts)

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
┌─────────────┐  ┌───────▼─────┐
│course_reg..│  │academic_years│
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
│course_off..│  │ semesters   │
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
        out: "/db/generated"
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
POST /academic/course-offering/:id/enroll - Enroll student in course offering
```

**Authentication**: Protected routes require `Authorization: Bearer <jwt-token>` header.
**Authorization**: Role-based access control enforced via middleware - enrollment endpoint restricted to students only (`RoleStudent = 3`).

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
**Status**: ✅ Fully implemented with course enrollment system

**Current Structure**:
```
modules/academic/
├── handlers/
│   └── course_enrollment.go    # Course enrollment endpoint
└── usecases/
    ├── course_enrollment.go     # Business logic & validation
    └── course_enrollment_test.go # Comprehensive unit tests
```

#### Implemented Features
- **Course Enrollment**: Students can enroll in course offerings
- **Business Validation**: 
  - Duplicate enrollment prevention
  - Capacity checking
  - Schedule conflict detection
- **Error Handling**: Detailed error responses for all failure scenarios
- **Test Coverage**: Full unit test suite with mock dependencies

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
**Current Status**: ✅ Testify framework configured and implemented

**Current Implementation**:
- **Unit Tests**: Complete test suite for course enrollment use case
- **Mocking**: Repository mocks using testify/mock
- **Test Suites**: Organized test suites with setup/teardown
- **Coverage**: Business logic validation, error scenarios, edge cases
- **Helper Function Tests**: Time calculations and overlap detection

**Test File**: `modules/academic/usecases/course_enrollment_test.go`

**Future Expansion**:
- Integration tests for repositories
- Handler tests with HTTP mocking
- Database tests with test containers

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
├── cmd/                  # Application entry point
│   └── main.go
├── config/               # Configuration management
│   └── config.go
├── common/               # Shared utilities
│   ├── base_response.go  # Standardized responses
│   └── validator.go      # Request validation
├── constants/            # System constants
│   └── constant.go       # Role definitions
├── middlewares/          # HTTP middleware
│   ├── jwt.go           # JWT authentication
│   └── access_control.go # Role-based access control
├── db/                   # Database layer
│   ├── generated/        # SQLC generated code
│   │   ├── models.go
│   │   ├── db.go
│   │   ├── users.sql.go
│   │   └── academic.sql.go
│   ├── migrations/       # Goose migration files
│   ├── repositories/     # Repository implementations
│   │   ├── users.go
│   │   └── academic.go
│   └── sql/              # SQL query definitions
│       ├── users.sql
│       └── academic.sql
├── modules/              # Feature modules
│   ├── auth/             # Authentication module
│   │   ├── handlers/
│   │   └── usecases/
│   └── academic/         # Academic management module
│       ├── handlers/
│       │   └── course_enrollment.go
│       └── usecases/
│           ├── course_enrollment.go
│           └── course_enrollment_test.go
├── docs/                 # Documentation
│   └── academic/
│       └── course-enrollment.md
├── config.json.example   # Configuration template
├── sqlc.yml             # SQLC configuration
├── go.mod               # Go module definition
├── ARCHITECTURE.md      # System architecture documentation
└── CLAUDE.md            # Development guidance
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
- **Middleware**: Authentication, CORS, rate limiting
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

#### Current Limitations
1. ✅ ~~**No Authentication Middleware**: JWT validation per endpoint~~ → **RESOLVED**: Centralized JWT middleware implemented
2. **Limited Error Types**: Generic error responses
3. **No API Versioning**: Single version assumption
4. **No Request Logging**: Limited observability
5. **No Health Checks**: Service monitoring gaps

#### Recommended Improvements
1. ✅ ~~**Add Middleware Layer**: Centralized cross-cutting concerns~~ → **RESOLVED**: JWT middleware implemented
2. **Enhance Error Handling**: Typed errors with proper HTTP codes  
3. **Add API Versioning**: Future compatibility
4. **Implement Observability**: Metrics, tracing, and monitoring
5. ✅ ~~**Add Testing Framework**: Comprehensive test coverage~~ → **RESOLVED**: Testify framework with comprehensive tests

---

*This document reflects the current architectural state of the SIAKAD system as of September 2025. The system has evolved from a proof-of-concept to a mature implementation with complete authentication, course enrollment functionality, and comprehensive testing. For development guidance and implementation patterns, refer to `CLAUDE.md`.*