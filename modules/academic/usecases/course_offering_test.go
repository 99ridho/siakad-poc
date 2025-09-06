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

// Mock repository for course offering tests
type MockCourseOfferingRepository struct {
	mock.Mock
}

func (m *MockCourseOfferingRepository) GetCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetCourse(ctx context.Context, id string) (generated.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.Course), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetCourseOfferingWithCourse(ctx context.Context, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetStudentEnrollmentsWithDetails(ctx context.Context, studentID string) ([]repositories.StudentEnrollmentWithDetails, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]repositories.StudentEnrollmentWithDetails), args.Error(1)
}

func (m *MockCourseOfferingRepository) CountCourseOfferingEnrollments(ctx context.Context, courseOfferingID string) (int64, error) {
	args := m.Called(ctx, courseOfferingID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCourseOfferingRepository) CheckEnrollmentExists(ctx context.Context, studentID, courseOfferingID string) (bool, error) {
	args := m.Called(ctx, studentID, courseOfferingID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCourseOfferingRepository) CreateEnrollment(ctx context.Context, studentID, courseOfferingID string) (generated.CourseRegistration, error) {
	args := m.Called(ctx, studentID, courseOfferingID)
	return args.Get(0).(generated.CourseRegistration), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetCourseOfferingsWithPagination(ctx context.Context, limit, offset int) ([]repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockCourseOfferingRepository) CountCourseOfferings(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCourseOfferingRepository) CreateCourseOffering(ctx context.Context, semesterID, courseID, sectionCode string, capacity int32, startTime time.Time) (generated.CourseOffering, error) {
	args := m.Called(ctx, semesterID, courseID, sectionCode, capacity, startTime)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockCourseOfferingRepository) UpdateCourseOffering(ctx context.Context, id, semesterID, courseID, sectionCode string, capacity int32, startTime time.Time) (generated.CourseOffering, error) {
	args := m.Called(ctx, id, semesterID, courseID, sectionCode, capacity, startTime)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockCourseOfferingRepository) DeleteCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(generated.CourseOffering), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetCourseOfferingByIDWithDetails(ctx context.Context, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

// Transaction-aware methods (required by interface)
func (m *MockCourseOfferingRepository) GetCourseOfferingWithCourseTx(txCtx *common.TxContext, id string) (repositories.CourseOfferingWithCourse, error) {
	args := m.Called(txCtx, id)
	return args.Get(0).(repositories.CourseOfferingWithCourse), args.Error(1)
}

func (m *MockCourseOfferingRepository) GetStudentEnrollmentsWithDetailsTx(txCtx *common.TxContext, studentID string) ([]repositories.StudentEnrollmentWithDetails, error) {
	args := m.Called(txCtx, studentID)
	return args.Get(0).([]repositories.StudentEnrollmentWithDetails), args.Error(1)
}

func (m *MockCourseOfferingRepository) CountCourseOfferingEnrollmentsTx(txCtx *common.TxContext, courseOfferingID string) (int64, error) {
	args := m.Called(txCtx, courseOfferingID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCourseOfferingRepository) CheckEnrollmentExistsTx(txCtx *common.TxContext, studentID, courseOfferingID string) (bool, error) {
	args := m.Called(txCtx, studentID, courseOfferingID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCourseOfferingRepository) CreateEnrollmentTx(txCtx *common.TxContext, studentID, courseOfferingID string) (generated.CourseRegistration, error) {
	args := m.Called(txCtx, studentID, courseOfferingID)
	return args.Get(0).(generated.CourseRegistration), args.Error(1)
}

// Test Suite
type CourseOfferingUseCaseTestSuite struct {
	suite.Suite
	useCase         *CourseOfferingUseCase
	mockRepo        *MockCourseOfferingRepository
	ctx             context.Context
	testTime        time.Time
	courseOfferUUID pgtype.UUID
	semesterUUID    pgtype.UUID
	courseUUID      pgtype.UUID
}

func (suite *CourseOfferingUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockCourseOfferingRepository)
	suite.useCase = NewCourseOfferingUseCase(suite.mockRepo)
	suite.ctx = context.Background()
	suite.testTime = time.Now()

	// Setup test UUIDs
	suite.courseOfferUUID = pgtype.UUID{
		Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Valid: true,
	}
	suite.semesterUUID = pgtype.UUID{
		Bytes: [16]byte{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		Valid: true,
	}
	suite.courseUUID = pgtype.UUID{
		Bytes: [16]byte{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
		Valid: true,
	}
}

func (suite *CourseOfferingUseCaseTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test successful pagination
func (suite *CourseOfferingUseCaseTestSuite) TestGetCourseOfferingsWithPagination_Success() {
	page := 1
	pageSize := 10
	limit := 10
	offset := 0
	totalRecords := int64(25)

	mockCourseOfferings := []repositories.CourseOfferingWithCourse{
		{
			CourseOfferingID:        suite.courseOfferUUID,
			SemesterID:              suite.semesterUUID,
			CourseID:                suite.courseUUID,
			SectionCode:             "A1",
			Capacity:                30,
			CourseOfferingStartTime: pgtype.Timestamptz{Time: suite.testTime, Valid: true},
			CourseCode:              "CS101",
			CourseName:              "Introduction to Computer Science",
			Credit:                  3,
		},
	}

	suite.mockRepo.On("GetCourseOfferingsWithPagination", suite.ctx, limit, offset).Return(mockCourseOfferings, nil)
	suite.mockRepo.On("CountCourseOfferings", suite.ctx).Return(totalRecords, nil)

	results, pagination, err := suite.useCase.GetCourseOfferingsWithPagination(suite.ctx, page, pageSize)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "CS101", results[0].CourseCode)
	assert.Equal(suite.T(), "Introduction to Computer Science", results[0].CourseName)
	assert.Equal(suite.T(), "A1", results[0].SectionCode)
	assert.Equal(suite.T(), int32(30), results[0].Capacity)
	assert.Equal(suite.T(), suite.testTime, results[0].StartTime)

	assert.NotNil(suite.T(), pagination)
	assert.Equal(suite.T(), page, pagination.Page)
	assert.Equal(suite.T(), pageSize, pagination.PageSize)
	assert.Equal(suite.T(), int(totalRecords), pagination.TotalRecords)
	assert.Equal(suite.T(), 3, pagination.TotalPages) // 25 / 10 = 3 pages
}

// Test pagination with default values
func (suite *CourseOfferingUseCaseTestSuite) TestGetCourseOfferingsWithPagination_DefaultValues() {
	page := 0     // Invalid, should default to 1
	pageSize := 0 // Invalid, should default to 10
	expectedLimit := 10
	expectedOffset := 0
	totalRecords := int64(5)

	suite.mockRepo.On("GetCourseOfferingsWithPagination", suite.ctx, expectedLimit, expectedOffset).Return([]repositories.CourseOfferingWithCourse{}, nil)
	suite.mockRepo.On("CountCourseOfferings", suite.ctx).Return(totalRecords, nil)

	results, pagination, err := suite.useCase.GetCourseOfferingsWithPagination(suite.ctx, page, pageSize)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), pagination)
	assert.Equal(suite.T(), 1, pagination.Page)      // Should default to 1
	assert.Equal(suite.T(), 10, pagination.PageSize) // Should default to 10
	assert.Empty(suite.T(), results)
}

// Test repository error during pagination
func (suite *CourseOfferingUseCaseTestSuite) TestGetCourseOfferingsWithPagination_RepositoryError() {
	page := 1
	pageSize := 10
	expectedError := errors.New("database connection error")

	suite.mockRepo.On("GetCourseOfferingsWithPagination", suite.ctx, pageSize, 0).Return([]repositories.CourseOfferingWithCourse{}, expectedError)

	results, pagination, err := suite.useCase.GetCourseOfferingsWithPagination(suite.ctx, page, pageSize)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Nil(suite.T(), results)
	assert.Nil(suite.T(), pagination)
}

// Test successful course offering creation
func (suite *CourseOfferingUseCaseTestSuite) TestCreateCourseOffering_Success() {
	req := CreateCourseOfferingRequest{
		CourseID:    "course-123",
		SemesterID:  "semester-456",
		SectionCode: "A1",
		Capacity:    30,
		StartTime:   suite.testTime,
	}

	expectedCourseOffering := generated.CourseOffering{
		ID: suite.courseOfferUUID,
	}

	suite.mockRepo.On("CreateCourseOffering", suite.ctx, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime).Return(expectedCourseOffering, nil)

	response, err := suite.useCase.CreateCourseOffering(suite.ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.ID)
}

// Test create course offering with repository error
func (suite *CourseOfferingUseCaseTestSuite) TestCreateCourseOffering_RepositoryError() {
	req := CreateCourseOfferingRequest{
		CourseID:    "course-123",
		SemesterID:  "semester-456",
		SectionCode: "A1",
		Capacity:    30,
		StartTime:   suite.testTime,
	}

	expectedError := errors.New("duplicate key violation")
	suite.mockRepo.On("CreateCourseOffering", suite.ctx, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime).Return(generated.CourseOffering{}, expectedError)

	response, err := suite.useCase.CreateCourseOffering(suite.ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), response.ID)
}

// Test successful course offering update
func (suite *CourseOfferingUseCaseTestSuite) TestUpdateCourseOffering_Success() {
	id := "course-offer-123"
	req := UpdateCourseOfferingRequest{
		CourseID:    "course-123",
		SemesterID:  "semester-456",
		SectionCode: "B2",
		Capacity:    25,
		StartTime:   suite.testTime,
	}

	expectedCourseOffering := generated.CourseOffering{
		ID: suite.courseOfferUUID,
	}

	suite.mockRepo.On("UpdateCourseOffering", suite.ctx, id, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime).Return(expectedCourseOffering, nil)

	response, err := suite.useCase.UpdateCourseOffering(suite.ctx, id, req)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.ID)
}

// Test update course offering not found
func (suite *CourseOfferingUseCaseTestSuite) TestUpdateCourseOffering_NotFound() {
	id := "course-offer-123"
	req := UpdateCourseOfferingRequest{
		CourseID:    "course-123",
		SemesterID:  "semester-456",
		SectionCode: "B2",
		Capacity:    25,
		StartTime:   suite.testTime,
	}

	suite.mockRepo.On("UpdateCourseOffering", suite.ctx, id, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime).Return(generated.CourseOffering{}, pgx.ErrNoRows)

	response, err := suite.useCase.UpdateCourseOffering(suite.ctx, id, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "course offering not found", err.Error())
	assert.Empty(suite.T(), response.ID)
}

// Test successful course offering deletion
func (suite *CourseOfferingUseCaseTestSuite) TestDeleteCourseOffering_Success() {
	id := "course-offer-123"

	expectedCourseOffering := generated.CourseOffering{
		ID: suite.courseOfferUUID,
	}

	suite.mockRepo.On("DeleteCourseOffering", suite.ctx, id).Return(expectedCourseOffering, nil)

	err := suite.useCase.DeleteCourseOffering(suite.ctx, id)

	assert.NoError(suite.T(), err)
}

// Test delete course offering not found
func (suite *CourseOfferingUseCaseTestSuite) TestDeleteCourseOffering_NotFound() {
	id := "course-offer-123"

	suite.mockRepo.On("DeleteCourseOffering", suite.ctx, id).Return(generated.CourseOffering{}, pgx.ErrNoRows)

	err := suite.useCase.DeleteCourseOffering(suite.ctx, id)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "course offering not found", err.Error())
}

// Test UUID to string conversion
func (suite *CourseOfferingUseCaseTestSuite) TestUuidToString() {
	uuid := pgtype.UUID{
		Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Valid: true,
	}

	result := uuidToString(uuid)
	expected := "01020304-0506-0708-090a-0b0c0d0e0f10"
	assert.Equal(suite.T(), expected, result)
}

// Test invalid UUID to string conversion
func (suite *CourseOfferingUseCaseTestSuite) TestUuidToString_Invalid() {
	uuid := pgtype.UUID{
		Valid: false,
	}

	result := uuidToString(uuid)
	assert.Equal(suite.T(), "", result)
}

// Run the test suite
func TestCourseOfferingUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CourseOfferingUseCaseTestSuite))
}
