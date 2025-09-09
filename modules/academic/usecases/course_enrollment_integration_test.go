// +build integration

package usecases

import (
	"context"
	"siakad-poc/common"
	"siakad-poc/db/generated"
	"siakad-poc/db/repositories"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Integration test suite for course enrollment with real database transactions
// To run these tests: go test -v -tags=integration ./modules/academic/usecases/
type EnrollmentIntegrationTestSuite struct {
	suite.Suite
	pool               *pgxpool.Pool
	useCase            *CourseEnrollmentUseCase
	repo               repositories.AcademicRepository
	txExecutor         common.TransactionExecutor
	ctx                context.Context
	testStudentID      string
	testCourseOfferingID string
	cleanup            []func() error
}

func (suite *EnrollmentIntegrationTestSuite) SetupSuite() {
	// Note: In a real integration test, you would set up a test database connection
	// For this example, we'll show the pattern but skip actual database setup
	suite.ctx = context.Background()
	
	// Example connection (would need real database in practice):
	// config := pgxpool.ParseConfig("postgres://test_user:test_pass@localhost:5432/test_siakad")
	// suite.pool, _ = pgxpool.ConnectConfig(suite.ctx, config)
	
	// For demonstration, we'll use mock setup
	suite.repo = repositories.NewDefaultAcademicRepository(suite.pool)
	suite.txExecutor = common.NewPgxTransactionExecutor(suite.pool)
	suite.useCase = NewCourseEnrollmentUseCase(suite.repo, suite.txExecutor)
	
	// Test data IDs (would be generated from test data setup)
	suite.testStudentID = "550e8400-e29b-41d4-a716-446655440001"
	suite.testCourseOfferingID = "550e8400-e29b-41d4-a716-446655440002"
}

func (suite *EnrollmentIntegrationTestSuite) TearDownSuite() {
	// Clean up test data
	for _, cleanupFunc := range suite.cleanup {
		cleanupFunc()
	}
	
	if suite.pool != nil {
		suite.pool.Close()
	}
}

// Test concurrent enrollment scenarios to verify transaction isolation
func (suite *EnrollmentIntegrationTestSuite) TestConcurrentEnrollment_LastSpotRace() {
	if suite.pool == nil {
		suite.T().Skip("Skipping integration test - no database connection")
		return
	}

	// Create a course offering with capacity 1 (only one spot available)
	// This would be set up in test data preparation
	
	numConcurrentStudents := 5
	studentIDs := make([]string, numConcurrentStudents)
	for i := 0; i < numConcurrentStudents; i++ {
		studentIDs[i] = generateTestStudentID(i)
	}

	var wg sync.WaitGroup
	results := make([]error, numConcurrentStudents)
	
	// Launch concurrent enrollment attempts
	for i := 0; i < numConcurrentStudents; i++ {
		wg.Add(1)
		go func(studentIndex int) {
			defer wg.Done()
			err := suite.useCase.EnrollStudent(suite.ctx, studentIDs[studentIndex], suite.testCourseOfferingID)
			results[studentIndex] = err
		}(i)
	}
	
	wg.Wait()
	
	// Verify that exactly one enrollment succeeded and others failed with capacity error
	successCount := 0
	capacityErrorCount := 0
	
	for _, err := range results {
		if err == nil {
			successCount++
		} else if IsEnrollmentError(err) {
			if errorType, ok := GetEnrollmentErrorType(err); ok && errorType == ErrCapacityExceeded {
				capacityErrorCount++
			}
		}
	}
	
	assert.Equal(suite.T(), 1, successCount, "Exactly one enrollment should succeed")
	assert.Equal(suite.T(), numConcurrentStudents-1, capacityErrorCount, "All other enrollments should fail with capacity error")
}

// Test transaction rollback behavior
func (suite *EnrollmentIntegrationTestSuite) TestTransactionRollback_DatabaseFailure() {
	if suite.pool == nil {
		suite.T().Skip("Skipping integration test - no database connection")
		return
	}

	// This test would simulate a database failure during enrollment
	// to ensure proper transaction rollback
	
	// Setup: Create test data that would cause a failure partway through enrollment
	// For example, invalid foreign key constraint that triggers after enrollment count
	
	err := suite.useCase.EnrollStudent(suite.ctx, suite.testStudentID, "invalid-course-offering-id")
	
	// Verify that no partial data was committed
	assert.Error(suite.T(), err)
	
	// Additional checks would verify that no enrollment record was created
	// and that capacity counts remain unchanged
}

// Test end-to-end enrollment flow with real database
func (suite *EnrollmentIntegrationTestSuite) TestEndToEndEnrollment_FullWorkflow() {
	if suite.pool == nil {
		suite.T().Skip("Skipping integration test - no database connection")
		return
	}

	// This test would cover the complete enrollment workflow:
	// 1. Create test academic year, semester, course, and course offering
	// 2. Create test student
	// 3. Perform enrollment
	// 4. Verify enrollment record is created correctly
	// 5. Verify capacity count is updated
	// 6. Test duplicate enrollment prevention
	// 7. Test schedule conflict detection with real data
	
	// Step 1: Setup test data (academic year, semester, course, course offering, student)
	testData := suite.setupTestData()
	defer suite.cleanupTestData(testData)
	
	// Step 2: Perform enrollment
	err := suite.useCase.EnrollStudent(suite.ctx, testData.StudentID, testData.CourseOfferingID)
	assert.NoError(suite.T(), err)
	
	// Step 3: Verify enrollment was created
	enrollmentExists, err := suite.repo.CheckEnrollmentExists(suite.ctx, testData.StudentID, testData.CourseOfferingID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), enrollmentExists)
	
	// Step 4: Verify capacity count
	currentCount, err := suite.repo.CountCourseOfferingEnrollments(suite.ctx, testData.CourseOfferingID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), currentCount)
	
	// Step 5: Test duplicate enrollment prevention
	err = suite.useCase.EnrollStudent(suite.ctx, testData.StudentID, testData.CourseOfferingID)
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	if errorType, ok := GetEnrollmentErrorType(err); ok {
		assert.Equal(suite.T(), ErrDuplicateEnrollment, errorType)
	}
}

// Test schedule conflict with multiple real enrollments
func (suite *EnrollmentIntegrationTestSuite) TestScheduleConflict_RealTimeData() {
	if suite.pool == nil {
		suite.T().Skip("Skipping integration test - no database connection")
		return
	}

	// Create overlapping course offerings with real time data
	testData := suite.setupScheduleConflictTestData()
	defer suite.cleanupTestData(testData)
	
	// Enroll in first course
	err := suite.useCase.EnrollStudent(suite.ctx, testData.StudentID, testData.FirstCourseOfferingID)
	assert.NoError(suite.T(), err)
	
	// Attempt to enroll in overlapping course
	err = suite.useCase.EnrollStudent(suite.ctx, testData.StudentID, testData.OverlappingCourseOfferingID)
	assert.Error(suite.T(), err)
	assert.True(suite.T(), IsEnrollmentError(err))
	if errorType, ok := GetEnrollmentErrorType(err); ok {
		assert.Equal(suite.T(), ErrScheduleConflict, errorType)
	}
	
	// Verify first enrollment is still valid
	exists, err := suite.repo.CheckEnrollmentExists(suite.ctx, testData.StudentID, testData.FirstCourseOfferingID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
	
	// Verify conflicting enrollment was not created
	exists, err = suite.repo.CheckEnrollmentExists(suite.ctx, testData.StudentID, testData.OverlappingCourseOfferingID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// Helper structures for test data
type TestData struct {
	StudentID        string
	CourseOfferingID string
	// Add other test data fields as needed
}

type ScheduleConflictTestData struct {
	StudentID                   string
	FirstCourseOfferingID       string
	OverlappingCourseOfferingID string
	NonOverlappingCourseOfferingID string
}

// Helper functions for test data setup and cleanup
func (suite *EnrollmentIntegrationTestSuite) setupTestData() *TestData {
	// In a real integration test, this would:
	// 1. Create academic year, semester, course, course offering records
	// 2. Create student record
	// 3. Return the IDs for use in tests
	
	return &TestData{
		StudentID:        generateTestStudentID(1),
		CourseOfferingID: generateTestCourseOfferingID(1),
	}
}

func (suite *EnrollmentIntegrationTestSuite) setupScheduleConflictTestData() *ScheduleConflictTestData {
	// Create course offerings with overlapping schedules
	// Course 1: 9:00-11:30 (3 credits)
	// Course 2: 10:00-12:00 (2 credits) - overlaps with Course 1
	// Course 3: 13:00-15:00 (2 credits) - no overlap
	
	return &ScheduleConflictTestData{
		StudentID:                      generateTestStudentID(1),
		FirstCourseOfferingID:          generateTestCourseOfferingID(1),
		OverlappingCourseOfferingID:    generateTestCourseOfferingID(2),
		NonOverlappingCourseOfferingID: generateTestCourseOfferingID(3),
	}
}

func (suite *EnrollmentIntegrationTestSuite) cleanupTestData(testData interface{}) {
	// Clean up test records from database
	// This would delete all test data created during the test
}

// Utility functions for generating test IDs
func generateTestStudentID(index int) string {
	// Generate unique student ID for testing
	return "test-student-" + string(rune('1'+index))
}

func generateTestCourseOfferingID(index int) string {
	// Generate unique course offering ID for testing
	return "test-course-offering-" + string(rune('1'+index))
}

// Benchmark test for enrollment performance
func (suite *EnrollmentIntegrationTestSuite) TestEnrollmentPerformance() {
	if suite.pool == nil {
		suite.T().Skip("Skipping integration test - no database connection")
		return
	}

	// Performance test to ensure enrollment operations complete within acceptable time
	numEnrollments := 100
	
	start := time.Now()
	
	for i := 0; i < numEnrollments; i++ {
		studentID := generateTestStudentID(i)
		courseOfferingID := generateTestCourseOfferingID(i % 10) // Distribute across 10 course offerings
		
		err := suite.useCase.EnrollStudent(suite.ctx, studentID, courseOfferingID)
		if err != nil && !IsBusinessRuleViolation(err) {
			// Only fail on non-business rule errors (system errors)
			suite.T().Fatalf("Unexpected error during enrollment %d: %v", i, err)
		}
	}
	
	duration := time.Since(start)
	
	// Assert performance criteria (adjust based on requirements)
	maxDuration := 10 * time.Second
	assert.True(suite.T(), duration < maxDuration, 
		"Enrollment performance test took %v, expected less than %v", duration, maxDuration)
	
	avgDuration := duration / time.Duration(numEnrollments)
	maxAvgDuration := 100 * time.Millisecond
	assert.True(suite.T(), avgDuration < maxAvgDuration,
		"Average enrollment time was %v, expected less than %v", avgDuration, maxAvgDuration)
}

// Run integration tests
func TestEnrollmentIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EnrollmentIntegrationTestSuite))
}