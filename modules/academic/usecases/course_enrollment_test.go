package usecases

import (
	"context"
	"errors"
	"siakad-poc/common"
	"siakad-poc/db/generated"
	"siakad-poc/db/repositories"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock repository for testing
type MockAcademicRepository struct {
	mock.Mock
}

func (m *MockAcademicRepository) GetCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockAcademicRepository) GetCourse(ctx context.Context, id string) (generated.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.Course), args.Error(1)
}

func (m *MockAcademicRepository) GetCourseOfferingWithCourse(ctx context.Context, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockAcademicRepository) GetStudentEnrollmentsWithDetails(ctx context.Context, studentID string) ([]repositories.StudentEnrollmentWithDetails, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]repositories.StudentEnrollmentWithDetails), args.Error(1)
}

func (m *MockAcademicRepository) CountCourseOfferingEnrollments(ctx context.Context, courseOfferingID string) (int64, error) {
	args := m.Called(ctx, courseOfferingID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAcademicRepository) CheckEnrollmentExists(ctx context.Context, studentID, courseOfferingID string) (bool, error) {
	args := m.Called(ctx, studentID, courseOfferingID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockAcademicRepository) CreateEnrollment(ctx context.Context, studentID, courseOfferingID string) (generated.CourseRegistration, error) {
	args := m.Called(ctx, studentID, courseOfferingID)
	return args.Get(0).(generated.CourseRegistration), args.Error(1)
}

// Course Offering CRUD methods (not used in enrollment tests, but required by interface)
func (m *MockAcademicRepository) GetCourseOfferingsWithPagination(ctx context.Context, limit, offset int) ([]repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockAcademicRepository) CountCourseOfferings(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAcademicRepository) CreateCourseOffering(ctx context.Context, semesterID, courseID, sectionCode string, capacity int32, startTime time.Time) (generated.CourseOffering, error) {
	args := m.Called(ctx, semesterID, courseID, sectionCode, capacity, startTime)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockAcademicRepository) UpdateCourseOffering(ctx context.Context, id, semesterID, courseID, sectionCode string, capacity int32, startTime time.Time) (generated.CourseOffering, error) {
	args := m.Called(ctx, id, semesterID, courseID, sectionCode, capacity, startTime)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockAcademicRepository) DeleteCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockAcademicRepository) GetCourseOfferingByIDWithDetails(ctx context.Context, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

// Transaction-aware methods (required by interface)
func (m *MockAcademicRepository) GetCourseOfferingWithCourseTx(txCtx *common.TxContext, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(txCtx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockAcademicRepository) GetStudentEnrollmentsWithDetailsTx(txCtx *common.TxContext, studentID string) ([]repositories.StudentEnrollmentWithDetails, error) {
	args := m.Called(txCtx, studentID)
	return args.Get(0).([]repositories.StudentEnrollmentWithDetails), args.Error(1)
}

func (m *MockAcademicRepository) CountCourseOfferingEnrollmentsTx(txCtx *common.TxContext, courseOfferingID string) (int64, error) {
	args := m.Called(txCtx, courseOfferingID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAcademicRepository) CheckEnrollmentExistsTx(txCtx *common.TxContext, studentID, courseOfferingID string) (bool, error) {
	args := m.Called(txCtx, studentID, courseOfferingID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockAcademicRepository) CreateEnrollmentTx(txCtx *common.TxContext, studentID, courseOfferingID string) (generated.CourseRegistration, error) {
	args := m.Called(txCtx, studentID, courseOfferingID)
	return args.Get(0).(generated.CourseRegistration), args.Error(1)
}

// Test Suite
type EnrollmentUseCaseTestSuite struct {
	suite.Suite
	useCase        *CourseEnrollmentUseCase
	mockRepo       *MockAcademicRepository
	mockTxExecutor *common.MockTransactionExecutor
	ctx            context.Context
	studentID      string
	courseID       string
}

func (suite *EnrollmentUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockAcademicRepository)
	suite.mockTxExecutor = new(common.MockTransactionExecutor)

	suite.useCase = &CourseEnrollmentUseCase{
		academicRepo: suite.mockRepo,
		txExecutor:   suite.mockTxExecutor,
	}

	suite.ctx = context.Background()
	suite.studentID = "550e8400-e29b-41d4-a716-446655440001"
	suite.courseID = "550e8400-e29b-41d4-a716-446655440002"
}

func (suite *EnrollmentUseCaseTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test successful enrollment
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_Success() {
	// Setup
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Mock expectations for transaction methods
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return([]repositories.StudentEnrollmentWithDetails{}, nil)
	suite.mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test duplicate enrollment
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_DuplicateEnrollment() {
	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(true, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrDuplicateEnrollment, errorType)
	assert.True(suite.T(), IsBusinessRuleViolation(err))
}

// Test course offering not found
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_CourseOfferingNotFound() {
	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(repositories.CourseOfferingWithCourse{}, pgx.ErrNoRows)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrCourseOfferingNotFound, errorType)
	assert.True(suite.T(), IsDataValidationError(err))
}

// Test capacity full
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_CapacityFull() {
	// Setup
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 10,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(10), nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrCapacityExceeded, errorType)
	assert.True(suite.T(), IsBusinessRuleViolation(err))
}

// Test schedule overlap
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_ScheduleOverlap() {
	// Setup - new course from 9:00-11:30 (3 credits * 50 min = 150 min)
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Existing enrollment from 10:00-11:40 (2 credits * 50 min = 100 min) - overlaps with new course
	existingEnrollments := []repositories.StudentEnrollmentWithDetails{
		{
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Credit: 2,
		},
	}

	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return(existingEnrollments, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrScheduleConflict, errorType)
	assert.True(suite.T(), IsBusinessRuleViolation(err))
}

// Test no schedule overlap
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_NoScheduleOverlap() {
	// Setup - new course from 9:00-11:30 (3 credits * 50 min = 150 min)
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Existing enrollment from 13:00-14:40 (2 credits * 50 min = 100 min) - no overlap
	existingEnrollments := []repositories.StudentEnrollmentWithDetails{
		{
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Credit: 2,
		},
	}

	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return(existingEnrollments, nil)
	suite.mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test repository error scenarios
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_RepositoryErrors() {
	// Test CheckEnrollmentExists error
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, errors.New("db error"))

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrDatabaseOperation, errorType)

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test GetCourseOfferingWithCourse error
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(repositories.CourseOfferingWithCourse{}, errors.New("db error"))

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok = GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrDatabaseOperation, errorType)
}

// Test data integrity validation
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_InvalidCourseOfferingData() {
	// Test invalid capacity (zero)
	courseOfferingWithInvalidCapacity := repositories.CourseOfferingWithCourse{
		Capacity: 0, // Invalid capacity
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithInvalidCapacity, nil)

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "capacity must be greater than 0")

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test invalid credit (zero)
	courseOfferingWithInvalidCredit := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 0, // Invalid credit
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithInvalidCredit, nil)

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "credit must be greater than 0")

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test invalid start time (NULL/invalid)
	courseOfferingWithInvalidStartTime := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Time{},
			Valid: false, // Invalid start time
		},
		Credit: 3,
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithInvalidStartTime, nil)

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "start time is not set")
}

// Test boundary conditions for capacity
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_CapacityBoundaryConditions() {
	// Test exactly at capacity (last spot available)
	courseOffering := repositories.CourseOfferingWithCourse{
		Capacity: 10,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Mock for exactly one spot left (9 enrolled out of 10)
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOffering, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(9), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return([]repositories.StudentEnrollmentWithDetails{}, nil)
	suite.mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.NoError(suite.T(), err)

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test exactly at full capacity (should fail)
	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOffering, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(10), nil)

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "full capacity (10/10)")
}

// Test schedule overlap edge cases
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_ScheduleOverlapEdgeCases() {
	// Test exact boundary case - courses that are adjacent but don't overlap
	// New course: 9:00-11:30 (3 credits = 150 minutes)
	// Existing course: 11:30-13:00 (2 credits = 100 minutes) - adjacent but no overlap
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	existingEnrollment := []repositories.StudentEnrollmentWithDetails{
		{
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 11, 30, 0, 0, time.UTC), // Starts exactly when first course ends
				Valid: true,
			},
			Credit: 2,
		},
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(5), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return(existingEnrollment, nil)
	suite.mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.NoError(suite.T(), err) // Should succeed - no overlap

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test 1-minute overlap case (should fail)
	// New course: 9:00-11:30 (3 credits = 150 minutes)
	// Existing course: 11:29-12:19 (1 credit = 50 minutes) - 1 minute overlap
	existingOverlapEnrollment := []repositories.StudentEnrollmentWithDetails{
		{
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 11, 29, 0, 0, time.UTC), // Starts 1 minute before first course ends
				Valid: true,
			},
			Credit: 1,
		},
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(5), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return(existingOverlapEnrollment, nil)

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	errorType, ok := GetEnrollmentErrorType(err)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), ErrScheduleConflict, errorType)
}

// Test handling of invalid existing enrollment data
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_InvalidExistingEnrollmentData() {
	courseOfferingWithCourse := repositories.CourseOfferingWithCourse{
		Capacity: 30,
		CourseOfferingStartTime: pgtype.Timestamptz{
			Time:  time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Credit: 3,
	}

	// Mix of valid and invalid existing enrollments - invalid ones should be skipped
	existingEnrollments := []repositories.StudentEnrollmentWithDetails{
		{
			// Invalid enrollment - no valid start time
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Time{},
				Valid: false,
			},
			Credit: 2,
		},
		{
			// Invalid enrollment - zero credits
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Credit: 0,
		},
		{
			// Valid enrollment - should not conflict
			CourseOfferingStartTime: pgtype.Timestamptz{
				Time:  time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Credit: 2,
		},
	}

	suite.mockRepo.On("CheckEnrollmentExistsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourseTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollmentsTx", mock.AnythingOfType("*common.TxContext"), suite.courseID).Return(int64(5), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetailsTx", mock.AnythingOfType("*common.TxContext"), suite.studentID).Return(existingEnrollments, nil)
	suite.mockRepo.On("CreateEnrollmentTx", mock.AnythingOfType("*common.TxContext"), suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.NoError(suite.T(), err) // Should succeed - invalid enrollments are skipped

}

// Helper function tests
func TestCalculateCourseEndTime(t *testing.T) {
	startTime := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)

	// Test 1 credit (50 minutes)
	endTime := calculateCourseEndTime(startTime, 1)
	expected := time.Date(2025, 1, 15, 9, 50, 0, 0, time.UTC)
	assert.Equal(t, expected, endTime)

	// Test 3 credits (150 minutes = 2.5 hours)
	endTime = calculateCourseEndTime(startTime, 3)
	expected = time.Date(2025, 1, 15, 11, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, endTime)

	// Test edge case: 0 credits (should return start time unchanged)
	endTime = calculateCourseEndTime(startTime, 0)
	assert.Equal(t, startTime, endTime)

	// Test edge case: negative credits (should return start time unchanged)
	endTime = calculateCourseEndTime(startTime, -1)
	assert.Equal(t, startTime, endTime)

	// Test large credit value (6 credits = 300 minutes = 5 hours)
	endTime = calculateCourseEndTime(startTime, 6)
	expected = time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, endTime)
}

func TestHasTimeOverlap(t *testing.T) {
	// Course 1: 9:00-11:00
	start1 := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)
	end1 := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)

	// Course 2: 10:00-12:00 (overlaps with Course 1)
	start2 := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	end2 := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test partial overlap
	assert.True(t, hasTimeOverlap(start1, end1, start2, end2))

	// Course 3: 11:00-13:00 (no overlap with Course 1 - adjacent)
	start3 := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)
	end3 := time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC)

	// Test no overlap (adjacent)
	assert.False(t, hasTimeOverlap(start1, end1, start3, end3))

	// Course 4: 8:00-9:00 (adjacent to Course 1, no overlap)
	start4 := time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC)
	end4 := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)

	// Test adjacent no overlap (before)
	assert.False(t, hasTimeOverlap(start1, end1, start4, end4))

	// Course 5: 9:30-10:30 (completely contained within Course 1)
	start5 := time.Date(2025, 1, 15, 9, 30, 0, 0, time.UTC)
	end5 := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)

	// Test complete containment
	assert.True(t, hasTimeOverlap(start1, end1, start5, end5))

	// Course 6: 8:00-12:00 (completely contains Course 1)
	start6 := time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC)
	end6 := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test being completely contained
	assert.True(t, hasTimeOverlap(start1, end1, start6, end6))

	// Course 7: 10:59-12:00 (1-minute overlap)
	start7 := time.Date(2025, 1, 15, 10, 59, 0, 0, time.UTC)
	end7 := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test 1-minute overlap
	assert.True(t, hasTimeOverlap(start1, end1, start7, end7))

	// Course 8: 12:00-14:00 (completely separate)
	start8 := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	end8 := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)

	// Test completely separate
	assert.False(t, hasTimeOverlap(start1, end1, start8, end8))

	// Course 9: 9:00-11:00 (exact same time)
	start9 := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)
	end9 := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)

	// Test exact same time range
	assert.True(t, hasTimeOverlap(start1, end1, start9, end9))
}

func TestConvertPgTimestamp(t *testing.T) {
	// Test valid timestamp
	validTime := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)
	pgTime := pgtype.Timestamptz{
		Time:  validTime,
		Valid: true,
	}

	result, err := convertPgTimestamp(pgTime)
	assert.NoError(t, err)
	assert.Equal(t, validTime, result)

	// Test invalid timestamp (Valid = false)
	invalidPgTime := pgtype.Timestamptz{
		Time:  time.Time{},
		Valid: false,
	}

	_, err = convertPgTimestamp(invalidPgTime)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid timestamp: database field is NULL or invalid")

	// Test edge case: valid timestamp with zero time
	zeroTimePg := pgtype.Timestamptz{
		Time:  time.Time{},
		Valid: true, // Valid but zero time
	}

	result, err = convertPgTimestamp(zeroTimePg)
	assert.NoError(t, err)
	assert.Equal(t, time.Time{}, result)

	// Test with timezone information
	timeWithTZ := time.Date(2025, 1, 15, 9, 0, 0, 0, time.FixedZone("UTC+7", 7*60*60))
	pgTimeWithTZ := pgtype.Timestamptz{
		Time:  timeWithTZ,
		Valid: true,
	}

	result, err = convertPgTimestamp(pgTimeWithTZ)
	assert.NoError(t, err)
	assert.Equal(t, timeWithTZ, result)
}

// Run the test suite
func TestEnrollmentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(EnrollmentUseCaseTestSuite))
}
