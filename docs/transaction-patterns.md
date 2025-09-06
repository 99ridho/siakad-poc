# Transaction Management Patterns in SIAKAD

## Overview

The SIAKAD system implements a consistent transaction management pattern to ensure ACID properties across multi-step database operations. This document explains the implementation and usage patterns.

## Architecture

### Transaction Utilities (`common/transaction.go`)

The system provides transaction management utilities that wrap pgx transaction handling:

```go
// Execute a function within a transaction
func withTransaction(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error

// Execute a function with shared transaction context
func withTxContext(ctx context.Context, pool *pgxpool.Pool, fn TxContextFunc) error
```

### Repository Pattern

Repositories implement both standard and transaction-aware methods:

**Standard Methods:** Use connection pool directly

```go
GetCourseOffering(ctx context.Context, id string) (CourseOffering, error)
```

**Transaction-aware Methods:** Accept transaction context

```go
GetCourseOfferingTx(txCtx *TxContext, id string) (CourseOffering, error)
```

### Use Case Layer

Use cases coordinate transactions through the TransactionExecutor interface:

```go
type CourseEnrollmentUseCase struct {
    academicRepo repositories.AcademicRepository
    txExecutor   common.TransactionExecutor
}

func NewCourseEnrollmentUseCase(academicRepo repositories.AcademicRepository, txExecutor common.TransactionExecutor) *CourseEnrollmentUseCase {
    return &CourseEnrollmentUseCase{
        academicRepo: academicRepo,
        txExecutor:   txExecutor,
    }
}

func (u *CourseEnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseOfferingID string) error {
    return u.txExecutor.WithTxContext(ctx, func(txCtx *common.TxContext) error {
        // All operations within this function share the same transaction
        return nil
    })
}
```

### TransactionExecutor Interface

The system uses dependency injection with the TransactionExecutor interface:

```go
// Interface for transaction execution
type TransactionExecutor interface {
    WithTxContext(ctx context.Context, fn func(*TxContext) error) error
}

// Production implementation
type PgxTransactionExecutor struct {
    pool *pgxpool.Pool
}

func (p *PgxTransactionExecutor) WithTxContext(ctx context.Context, fn func(*TxContext) error) error {
    return withTxContext(ctx, p.pool, fn)
}
```

## Implementation Example: Course Enrollment

The course enrollment process demonstrates transaction usage:

### Problem Without Transactions

1. Check enrollment exists → `FALSE`
2. Count current enrollments → `9/10`
3. **[CONCURRENT ENROLLMENT OCCURS]**
4. Create enrollment → `11/10` (EXCEEDS CAPACITY)

### Solution With Transactions

```go
func (u *CourseEnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseOfferingID string) error {
    return u.txExecutor.WithTxContext(ctx, func(txCtx *common.TxContext) error {
        // 1. Check existing enrollment (within transaction)
        exists, err := u.academicRepo.CheckEnrollmentExistsTx(txCtx, studentID, courseOfferingID)
        if exists || err != nil {
            return handleError(err, "enrollment check failed")
        }

        // 2. Get course details (consistent read)
        courseOffering, err := u.academicRepo.GetCourseOfferingWithCourseTx(txCtx, courseOfferingID)
        if err != nil {
            return handleError(err, "course offering not found")
        }

        // 3. Check capacity (consistent count within transaction)
        count, err := u.academicRepo.CountCourseOfferingEnrollmentsTx(txCtx, courseOfferingID)
        if err != nil || count >= int64(courseOffering.Capacity) {
            return handleError(err, "capacity exceeded")
        }

        // 4. Validate schedule conflicts (consistent student data)
        enrollments, err := u.academicRepo.GetStudentEnrollmentsWithDetailsTx(txCtx, studentID)
        if err != nil {
            return handleError(err, "failed to get student enrollments")
        }

        if hasScheduleConflict(courseOffering, enrollments) {
            return errors.New("schedule conflict detected")
        }

        // 5. Create enrollment (atomic operation)
        _, err = u.academicRepo.CreateEnrollmentTx(txCtx, studentID, courseOfferingID)
        return handleError(err, "enrollment creation failed")
    })
}
```

## Benefits

### 1. **Data Consistency**

- All reads within transaction see consistent snapshot
- Prevents race conditions between validation and creation

### 2. **Atomicity**

- Either all operations succeed or all fail
- No partial state changes

### 3. **Isolation**

- Concurrent transactions don't interfere
- Capacity limits enforced correctly under load

### 4. **Rollback Safety**

- Automatic rollback on any error
- Database remains in consistent state

## Testing Patterns

### Unit Testing with Mocks

```go
// MockTransactionExecutor bypasses actual transactions for unit testing
type MockTransactionExecutor struct {
    mock.Mock
}

func (m *MockTransactionExecutor) WithTxContext(ctx context.Context, fn func(*common.TxContext) error) error {
    // Execute function directly without transaction overhead
    mockTxCtx := &common.TxContext{} // Minimal mock context
    return fn(mockTxCtx)
}

func TestTransactionRollback(t *testing.T) {
    mockRepo := &MockAcademicRepository{}
    mockTxExecutor := &MockTransactionExecutor{}

    // Mock successful validations
    mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), studentID, courseID).Return(false, nil)
    mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), courseID).Return(validCourse, nil)
    mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), courseID).Return(int64(5), nil)
    mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), studentID).Return([]Enrollment{}, nil)

    // Mock enrollment failure
    mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), studentID, courseID).Return(nil, errors.New("db error"))

    useCase := NewCourseEnrollmentUseCase(mockRepo, mockTxExecutor)
    err := useCase.EnrollStudent(ctx, studentID, courseID)

    assert.Error(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Integration Testing

Integration tests should use test databases to verify actual transaction behavior:

```go
func TestActualTransactionRollback(t *testing.T) {
    testDB := setupTestDatabase(t)
    defer cleanupTestDatabase(testDB)

    // Create use case with real repository and transaction executor
    repo := repositories.NewDefaultAcademicRepository(testDB)
    txExecutor := common.NewPgxTransactionExecutor(testDB)
    useCase := usecases.NewCourseEnrollmentUseCase(repo, txExecutor)

    // Setup test data
    student := createTestStudent(testDB)
    course := createTestCourseOffering(testDB)

    // Force failure after validations
    forceEnrollmentTableConstraintViolation(testDB)

    // Attempt enrollment
    err := useCase.EnrollStudent(ctx, student.ID, course.ID)
    assert.Error(t, err)

    // Verify rollback - no enrollment should exist
    count := countEnrollments(testDB, student.ID, course.ID)
    assert.Equal(t, 0, count)
}
```

## Best Practices

### 1. **Use Case Transaction Boundary**

- Start transactions at use case level through TransactionExecutor interface
- Keep transaction scope focused on business operations
- Use dependency injection for testability

### 2. **Repository Dual Interface**

- Provide both standard and transaction-aware methods
- Use `Tx` suffix for transaction methods
- Transaction methods accept `*common.TxContext` parameter

### 3. **Dependency Injection Pattern**

- Inject TransactionExecutor interface into use cases
- Use PgxTransactionExecutor for production
- Use MockTransactionExecutor for unit tests

### 4. **Error Handling**

- Use structured error wrapping
- Provide context for transaction failures
- Automatic rollback on any error within transaction

### 5. **Testing Strategy**

- Unit tests with MockTransactionExecutor and repository mocks
- Integration tests with real PgxTransactionExecutor
- Load testing for concurrent transaction scenarios

## Migration Guide

### For Existing Use Cases

1. **Add TransactionExecutor dependency:**

```go
type UseCase struct {
    repo       Repository
    txExecutor common.TransactionExecutor  // Add this interface
}

func NewUseCase(repo Repository, txExecutor common.TransactionExecutor) *UseCase {
    return &UseCase{
        repo:       repo,
        txExecutor: txExecutor,
    }
}
```

2. **Wrap critical operations:**

```go
func (u *UseCase) CriticalOperation(ctx context.Context, params) error {
    return u.txExecutor.WithTxContext(ctx, func(txCtx *common.TxContext) error {
        // Convert repo calls to Tx variants
        return u.repo.CreateSomethingTx(txCtx, params)
    })
}
```

3. **Update module wiring:**

```go
func NewModule(pool *pgxpool.Pool) *Module {
    txExecutor := common.NewPgxTransactionExecutor(pool)
    repo := repositories.NewRepository(pool)
    useCase := NewUseCase(repo, txExecutor)
    return &Module{useCase: useCase}
}
```

3. **Update repositories:**

- Add transaction-aware methods to interfaces
- Implement using `generated.Queries.WithTx()`

### For New Features

- Always consider if operations need transaction protection
- Use transaction-aware methods for multi-step operations
- Implement comprehensive test coverage

## Performance Considerations

- Transactions have overhead - use judiciously
- Keep transaction scope minimal
- Consider connection pool sizing for transaction load
- Monitor for deadlocks in complex transaction scenarios

## Monitoring

- Log transaction begin/commit/rollback events
- Monitor transaction duration
- Alert on high rollback rates
- Track concurrent transaction conflicts
