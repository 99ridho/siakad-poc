package usecases

import (
	"context"
	"fmt"
	"siakad-poc/common"
	"siakad-poc/db/repositories"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

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

// EnrollStudent enrolls a student in a course offering after validating business rules.
// Business Rules Validated:
// 1. No duplicate enrollment - student cannot enroll twice in the same course offering
// 2. Capacity check - enrollment count must be less than course offering capacity
// 3. Schedule conflict detection - new course cannot overlap with existing enrollments
//    - Each credit = 50 minutes of class time
//    - Schedule overlap is calculated based on start_time + (credit * 50 minutes)
func (u *CourseEnrollmentUseCase) EnrollStudent(ctx context.Context, studentID, courseOfferingID string) error {
	// Execute all enrollment operations within a transaction to ensure ACID properties
	// This prevents race conditions and ensures data consistency across all validation steps
	return u.txExecutor.WithTxContext(ctx, func(txCtx *common.TxContext) error {
		// Business Rule 1: No Enrollment Duplication
		// Check if student is already enrolled in this course offering (with transaction)
		exists, err := u.academicRepo.CheckEnrollmentExistsTx(txCtx, studentID, courseOfferingID)
		if err != nil {
			return NewDatabaseOperationError("check enrollment existence", err)
		}
		if exists {
			return NewDuplicateEnrollmentError(studentID, courseOfferingID)
		}

		// Retrieve course offering with course details (with transaction for consistent read)
		// This ensures we get the latest data within the transaction context
		courseOfferingWithCourse, err := u.academicRepo.GetCourseOfferingWithCourseTx(txCtx, courseOfferingID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return NewCourseOfferingNotFoundError(courseOfferingID)
			}
			return NewDatabaseOperationError("get course offering details", err)
		}

		// Validate course offering data integrity
		if courseOfferingWithCourse.Capacity <= 0 {
			return NewInvalidCourseDataError("capacity", "must be greater than 0")
		}
		if courseOfferingWithCourse.Credit <= 0 {
			return NewInvalidCourseDataError("credit", "must be greater than 0")
		}
		if !courseOfferingWithCourse.CourseOfferingStartTime.Valid {
			return NewInvalidCourseDataError("start time", "is not set")
		}

		// Business Rule 2: Capacity Validation
		// Check capacity - ensure enrollment count is less than capacity (with transaction for consistent read)
		currentEnrollmentCount, err := u.academicRepo.CountCourseOfferingEnrollmentsTx(txCtx, courseOfferingID)
		if err != nil {
			return NewDatabaseOperationError("count current enrollments", err)
		}
		if currentEnrollmentCount >= int64(courseOfferingWithCourse.Capacity) {
			return NewCapacityExceededError(currentEnrollmentCount, int64(courseOfferingWithCourse.Capacity))
		}

		// Business Rule 3: Schedule Conflict Detection
		// Check for schedule overlaps with student's existing enrollments (with transaction)
		existingEnrollments, err := u.academicRepo.GetStudentEnrollmentsWithDetailsTx(txCtx, studentID)
		if err != nil {
			return NewDatabaseOperationError("get student's existing enrollments", err)
		}

		// Calculate the time range for the new course offering
		// Formula: end_time = start_time + (credit * 50 minutes)
		newCourseStartTime, err := convertPgTimestamp(courseOfferingWithCourse.CourseOfferingStartTime)
		if err != nil {
			return NewInvalidTimestampError("new course start time")
		}
		newCourseEndTime := calculateCourseEndTime(newCourseStartTime, courseOfferingWithCourse.Credit)

		// Validate against all existing enrollments for schedule conflicts
		for _, enrollment := range existingEnrollments {
			// Skip invalid enrollment data
			if !enrollment.CourseOfferingStartTime.Valid || enrollment.Credit <= 0 {
				continue
			}

			existingStartTime, err := convertPgTimestamp(enrollment.CourseOfferingStartTime)
			if err != nil {
				return NewInvalidTimestampError("existing course start time")
			}
			existingEndTime := calculateCourseEndTime(existingStartTime, enrollment.Credit)

			// Check for time overlap using inclusive boundary logic
			if hasTimeOverlap(newCourseStartTime, newCourseEndTime, existingStartTime, existingEndTime) {
				newCourseTime := fmt.Sprintf("%s-%s", newCourseStartTime.Format("15:04"), newCourseEndTime.Format("15:04"))
				existingCourseTime := fmt.Sprintf("%s-%s", existingStartTime.Format("15:04"), existingEndTime.Format("15:04"))
				return NewScheduleConflictError(newCourseTime, existingCourseTime)
			}
		}

		// All business rules validated successfully - create the enrollment
		// This operation is within the transaction to ensure atomic behavior
		_, err = u.academicRepo.CreateEnrollmentTx(txCtx, studentID, courseOfferingID)
		if err != nil {
			return NewDatabaseOperationError("create enrollment", err)
		}

		return nil
	})
}

// calculateCourseEndTime calculates the end time of a course based on its start time and credit hours.
// Business Rule: Each credit hour equals 50 minutes of class time.
// Formula: end_time = start_time + (credits * 50 minutes)
// Example: 3-credit course starting at 9:00 AM ends at 11:30 AM (9:00 + 150 minutes)
func calculateCourseEndTime(startTime time.Time, credits int32) time.Time {
	if credits <= 0 {
		return startTime // No duration if invalid credits
	}
	durationMinutes := int(credits) * 50
	return startTime.Add(time.Duration(durationMinutes) * time.Minute)
}

// hasTimeOverlap checks if two time ranges overlap using inclusive boundary logic.
// Two time ranges overlap if: start1 < end2 AND start2 < end1
// This handles all overlap scenarios including:
// - Partial overlaps (start1 < start2 < end1 < end2)  
// - Complete containment (start1 <= start2 && end2 <= end1)
// - Adjacent ranges are NOT considered overlapping (end1 == start2)
// Example: [9:00-11:00] and [10:00-12:00] overlap, but [9:00-11:00] and [11:00-13:00] do not
func hasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

// convertPgTimestamp safely converts pgtype.Timestamptz to standard time.Time.
// Returns error if the PostgreSQL timestamp is marked as invalid/NULL.
// This prevents runtime panics when working with potentially NULL database fields.
func convertPgTimestamp(pgTime pgtype.Timestamptz) (time.Time, error) {
	if !pgTime.Valid {
		return time.Time{}, NewInvalidTimestampError("database field is NULL or invalid")
	}
	return pgTime.Time, nil
}
