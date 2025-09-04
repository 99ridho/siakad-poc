package usecases

import (
	"context"
	"errors"
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

// Test Suite
type EnrollmentUseCaseTestSuite struct {
	suite.Suite
	useCase   *CourseEnrollmentUseCase
	mockRepo  *MockAcademicRepository
	ctx       context.Context
	studentID string
	courseID  string
}

func (suite *EnrollmentUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockAcademicRepository)
	suite.useCase = NewCourseEnrollmentUseCase(suite.mockRepo)
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

	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollments", suite.ctx, suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetails", suite.ctx, suite.studentID).Return([]repositories.StudentEnrollmentWithDetails{}, nil)
	suite.mockRepo.On("CreateEnrollment", suite.ctx, suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test duplicate enrollment
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_DuplicateEnrollment() {
	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(true, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "already enrolled")
}

// Test course offering not found
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_CourseOfferingNotFound() {
	// Mock expectations
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(repositories.CourseOfferingWithCourse{}, pgx.ErrNoRows)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "course offering not found")
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
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollments", suite.ctx, suite.courseID).Return(int64(10), nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "full capacity")
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
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollments", suite.ctx, suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetails", suite.ctx, suite.studentID).Return(existingEnrollments, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "schedule conflict")
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
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(courseOfferingWithCourse, nil)
	suite.mockRepo.On("CountCourseOfferingEnrollments", suite.ctx, suite.courseID).Return(int64(10), nil)
	suite.mockRepo.On("GetStudentEnrollmentsWithDetails", suite.ctx, suite.studentID).Return(existingEnrollments, nil)
	suite.mockRepo.On("CreateEnrollment", suite.ctx, suite.studentID, suite.courseID).Return(generated.CourseRegistration{}, nil)

	// Execute
	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test repository error scenarios
func (suite *EnrollmentUseCaseTestSuite) TestEnrollStudent_RepositoryErrors() {
	// Test CheckEnrollmentExists error
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, errors.New("db error"))

	err := suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to check enrollment existence")

	// Reset mock for next test
	suite.mockRepo.ExpectedCalls = nil
	suite.mockRepo.Calls = nil

	// Test GetCourseOfferingWithCourse error
	suite.mockRepo.On("CheckEnrollmentExists", suite.ctx, suite.studentID, suite.courseID).Return(false, nil)
	suite.mockRepo.On("GetCourseOfferingWithCourse", suite.ctx, suite.courseID).Return(repositories.CourseOfferingWithCourse{}, errors.New("db error"))

	err = suite.useCase.EnrollStudent(suite.ctx, suite.studentID, suite.courseID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get course offering details")
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
}

func TestHasTimeOverlap(t *testing.T) {
	// Course 1: 9:00-11:00
	start1 := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)
	end1 := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)

	// Course 2: 10:00-12:00 (overlaps with Course 1)
	start2 := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	end2 := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Test overlap
	assert.True(t, hasTimeOverlap(start1, end1, start2, end2))

	// Course 3: 11:00-13:00 (no overlap with Course 1)
	start3 := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)
	end3 := time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC)

	// Test no overlap
	assert.False(t, hasTimeOverlap(start1, end1, start3, end3))

	// Course 4: 8:00-9:00 (adjacent to Course 1, no overlap)
	start4 := time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC)
	end4 := time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC)

	// Test adjacent no overlap
	assert.False(t, hasTimeOverlap(start1, end1, start4, end4))
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

	// Test invalid timestamp
	invalidPgTime := pgtype.Timestamptz{
		Time:  time.Time{},
		Valid: false,
	}

	_, err = convertPgTimestamp(invalidPgTime)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid timestamp")
}

// Run the test suite
func TestEnrollmentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(EnrollmentUseCaseTestSuite))
}
