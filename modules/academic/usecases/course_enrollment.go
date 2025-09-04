package usecases

import (
	"context"
	"fmt"
	"siakad-poc/db/repositories"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

type CourseEnrollmentUseCase struct {
	academicRepo repositories.AcademicRepository
}

func NewCourseEnrollmentUseCase(academicRepo repositories.AcademicRepository) *CourseEnrollmentUseCase {
	return &CourseEnrollmentUseCase{
		academicRepo: academicRepo,
	}
}

func (u *CourseEnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseOfferingID string) error {
	// 1. Check if student is already enrolled in this course offering
	exists, err := u.academicRepo.CheckEnrollmentExists(ctx, studentID, courseOfferingID)
	if err != nil {
		return errors.Wrap(err, "failed to check enrollment existence")
	}
	if exists {
		return errors.New("student is already enrolled in this course offering")
	}

	// 2. Get course offering with course details (needed for capacity and schedule validation)
	courseOfferingWithCourse, err := u.academicRepo.GetCourseOfferingWithCourse(ctx, courseOfferingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("course offering not found")
		}
		return errors.Wrap(err, "failed to get course offering details")
	}

	// 3. Check capacity - ensure enrollment count is less than capacity
	currentEnrollmentCount, err := u.academicRepo.CountCourseOfferingEnrollments(ctx, courseOfferingID)
	if err != nil {
		return errors.Wrap(err, "failed to count current enrollments")
	}
	if currentEnrollmentCount >= int64(courseOfferingWithCourse.Capacity) {
		return errors.New("course offering is at full capacity")
	}

	// 4. Check for schedule overlaps with student's existing enrollments
	existingEnrollments, err := u.academicRepo.GetStudentEnrollmentsWithDetails(ctx, studentID)
	if err != nil {
		return errors.Wrap(err, "failed to get student's existing enrollments")
	}

	// Calculate the time range for the new course offering
	newCourseStartTime, err := convertPgTimestamp(courseOfferingWithCourse.CourseOfferingStartTime)
	if err != nil {
		return errors.Wrap(err, "failed to parse new course start time")
	}
	newCourseEndTime := calculateCourseEndTime(newCourseStartTime, courseOfferingWithCourse.Credit)

	// Check for overlaps with existing enrollments
	for _, enrollment := range existingEnrollments {
		existingStartTime, err := convertPgTimestamp(enrollment.CourseOfferingStartTime)
		if err != nil {
			return errors.Wrap(err, "failed to parse existing course start time")
		}
		existingEndTime := calculateCourseEndTime(existingStartTime, enrollment.Credit)

		if hasTimeOverlap(newCourseStartTime, newCourseEndTime, existingStartTime, existingEndTime) {
			return fmt.Errorf("schedule conflict: new course overlaps with existing enrollment")
		}
	}

	// 5. Create the enrollment if all validations pass
	_, err = u.academicRepo.CreateEnrollment(ctx, studentID, courseOfferingID)
	if err != nil {
		return errors.Wrap(err, "failed to create enrollment")
	}

	return nil
}

// Helper function to calculate course end time based on credits
// Each credit = 50 minutes
func calculateCourseEndTime(startTime time.Time, credits int32) time.Time {
	durationMinutes := int(credits) * 50
	return startTime.Add(time.Duration(durationMinutes) * time.Minute)
}

// Helper function to check if two time ranges overlap
func hasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	// Two ranges overlap if start1 < end2 AND start2 < end1
	return start1.Before(end2) && start2.Before(end1)
}

// Helper function to convert pgtype.Timestamptz to time.Time
func convertPgTimestamp(pgTime pgtype.Timestamptz) (time.Time, error) {
	if !pgTime.Valid {
		return time.Time{}, errors.New("invalid timestamp")
	}
	return pgTime.Time, nil
}
