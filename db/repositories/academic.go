package repositories

import (
	"context"
	"errors"
	"siakad-poc/db/generated"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseOfferingWithCourse struct {
	CourseOfferingID        pgtype.UUID
	SemesterID              pgtype.UUID
	CourseID                pgtype.UUID
	SectionCode             string
	Capacity                int32
	CourseOfferingStartTime pgtype.Timestamptz
	CourseCode              string
	CourseName              string
	Credit                  int32
}

type StudentEnrollmentWithDetails struct {
	RegistrationID          pgtype.UUID
	StudentID               pgtype.UUID
	CourseOfferingID        pgtype.UUID
	RegistrationCreatedAt   pgtype.Timestamptz
	CourseOfferingStartTime pgtype.Timestamptz
	Credit                  int32
}

type AcademicRepository interface {
	GetCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error)
	GetCourse(ctx context.Context, id string) (generated.Course, error)
	GetCourseOfferingWithCourse(ctx context.Context, id string) (CourseOfferingWithCourse, error)
	GetStudentEnrollmentsWithDetails(ctx context.Context, studentID string) ([]StudentEnrollmentWithDetails, error)
	CountCourseOfferingEnrollments(ctx context.Context, courseOfferingID string) (int64, error)
	CheckEnrollmentExists(ctx context.Context, studentID, courseOfferingID string) (bool, error)
	CreateEnrollment(ctx context.Context, studentID, courseOfferingID string) (generated.CourseRegistration, error)
}

type DefaultAcademicRepository struct {
	query *generated.Queries
	pool  *pgxpool.Pool
}

func NewDefaultAcademicRepository(pool *pgxpool.Pool) *DefaultAcademicRepository {
	return &DefaultAcademicRepository{
		query: generated.New(pool),
		pool:  pool,
	}
}

func (r *DefaultAcademicRepository) GetCourseOffering(ctx context.Context, id string) (generated.CourseOffering, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(id)
	if err != nil {
		return generated.CourseOffering{}, errors.New("can't parse course offering id as uuid")
	}

	return r.query.GetCourseOffering(ctx, uuidID)
}

func (r *DefaultAcademicRepository) GetCourse(ctx context.Context, id string) (generated.Course, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(id)
	if err != nil {
		return generated.Course{}, errors.New("can't parse course id as uuid")
	}

	return r.query.GetCourse(ctx, uuidID)
}

func (r *DefaultAcademicRepository) GetCourseOfferingWithCourse(ctx context.Context, id string) (CourseOfferingWithCourse, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(id)
	if err != nil {
		return CourseOfferingWithCourse{}, errors.New("can't parse course offering id as uuid")
	}

	row, err := r.query.GetCourseOfferingWithCourse(ctx, uuidID)
	if err != nil {
		return CourseOfferingWithCourse{}, err
	}

	return CourseOfferingWithCourse{
		CourseOfferingID:        row.CourseOfferingID,
		SemesterID:              row.SemesterID,
		CourseID:                row.CourseID,
		SectionCode:             row.SectionCode,
		Capacity:                row.Capacity,
		CourseOfferingStartTime: row.CourseOfferingStartTime,
		CourseCode:              row.CourseCode,
		CourseName:              row.CourseName,
		Credit:                  row.Credit,
	}, nil
}

func (r *DefaultAcademicRepository) GetStudentEnrollmentsWithDetails(ctx context.Context, studentID string) ([]StudentEnrollmentWithDetails, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(studentID)
	if err != nil {
		return nil, errors.New("can't parse student id as uuid")
	}

	rows, err := r.query.GetStudentEnrollmentsWithDetails(ctx, uuidID)
	if err != nil {
		return nil, err
	}

	var enrollments []StudentEnrollmentWithDetails
	for _, row := range rows {
		enrollments = append(enrollments, StudentEnrollmentWithDetails{
			RegistrationID:          row.RegistrationID,
			StudentID:               row.StudentID,
			CourseOfferingID:        row.CourseOfferingID,
			RegistrationCreatedAt:   row.RegistrationCreatedAt,
			CourseOfferingStartTime: row.CourseOfferingStartTime,
			Credit:                  row.Credit,
		})
	}

	return enrollments, nil
}

func (r *DefaultAcademicRepository) CountCourseOfferingEnrollments(ctx context.Context, courseOfferingID string) (int64, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(courseOfferingID)
	if err != nil {
		return 0, errors.New("can't parse course offering id as uuid")
	}

	return r.query.CountCourseOfferingEnrollments(ctx, uuidID)
}

func (r *DefaultAcademicRepository) CheckEnrollmentExists(ctx context.Context, studentID, courseOfferingID string) (bool, error) {
	var studentUUID, courseOfferingUUID pgtype.UUID
	err := studentUUID.Scan(studentID)
	if err != nil {
		return false, errors.New("can't parse student id as uuid")
	}
	err = courseOfferingUUID.Scan(courseOfferingID)
	if err != nil {
		return false, errors.New("can't parse course offering id as uuid")
	}

	params := generated.CheckEnrollmentExistsParams{
		StudentID:        studentUUID,
		CourseOfferingID: courseOfferingUUID,
	}

	return r.query.CheckEnrollmentExists(ctx, params)
}

func (r *DefaultAcademicRepository) CreateEnrollment(ctx context.Context, studentID, courseOfferingID string) (generated.CourseRegistration, error) {
	var studentUUID, courseOfferingUUID pgtype.UUID
	err := studentUUID.Scan(studentID)
	if err != nil {
		return generated.CourseRegistration{}, errors.New("can't parse student id as uuid")
	}
	err = courseOfferingUUID.Scan(courseOfferingID)
	if err != nil {
		return generated.CourseRegistration{}, errors.New("can't parse course offering id as uuid")
	}

	params := generated.CreateEnrollmentParams{
		StudentID:        studentUUID,
		CourseOfferingID: courseOfferingUUID,
	}

	return r.query.CreateEnrollment(ctx, params)
}

// Helper function to convert pgtype.Timestamptz to time.Time
func (r *DefaultAcademicRepository) ConvertPgTimestamp(pgTime pgtype.Timestamptz) (time.Time, error) {
	if !pgTime.Valid {
		return time.Time{}, errors.New("invalid timestamp")
	}
	return pgTime.Time, nil
}