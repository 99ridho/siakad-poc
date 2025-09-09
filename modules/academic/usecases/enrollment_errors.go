package usecases

import "fmt"

// EnrollmentError represents domain-specific errors in the course enrollment process
type EnrollmentError struct {
	Type    EnrollmentErrorType
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *EnrollmentError) Error() string {
	return e.Message
}

// EnrollmentErrorType defines the different types of enrollment errors
type EnrollmentErrorType string

const (
	// Business rule violations
	ErrDuplicateEnrollment      EnrollmentErrorType = "DUPLICATE_ENROLLMENT"
	ErrCapacityExceeded         EnrollmentErrorType = "CAPACITY_EXCEEDED"
	ErrScheduleConflict         EnrollmentErrorType = "SCHEDULE_CONFLICT"
	
	// Data validation errors
	ErrCourseOfferingNotFound   EnrollmentErrorType = "COURSE_OFFERING_NOT_FOUND"
	ErrInvalidCourseData        EnrollmentErrorType = "INVALID_COURSE_DATA"
	ErrInvalidTimestamp         EnrollmentErrorType = "INVALID_TIMESTAMP"
	
	// System errors
	ErrDatabaseOperation        EnrollmentErrorType = "DATABASE_OPERATION"
	ErrTransactionFailed        EnrollmentErrorType = "TRANSACTION_FAILED"
)

// NewDuplicateEnrollmentError creates an error for duplicate enrollment attempts
func NewDuplicateEnrollmentError(studentID, courseOfferingID string) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrDuplicateEnrollment,
		Message: "Student is already enrolled in this course offering",
		Details: map[string]interface{}{
			"student_id":         studentID,
			"course_offering_id": courseOfferingID,
		},
	}
}

// NewCapacityExceededError creates an error for capacity violations
func NewCapacityExceededError(currentCount, maxCapacity int64) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrCapacityExceeded,
		Message: fmt.Sprintf("Course offering is at full capacity (%d/%d)", currentCount, maxCapacity),
		Details: map[string]interface{}{
			"current_enrollment": currentCount,
			"max_capacity":       maxCapacity,
		},
	}
}

// NewScheduleConflictError creates an error for schedule conflicts
func NewScheduleConflictError(newCourseTime, existingCourseTime string) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrScheduleConflict,
		Message: fmt.Sprintf("Schedule conflict: new course (%s) overlaps with existing enrollment (%s)", newCourseTime, existingCourseTime),
		Details: map[string]interface{}{
			"new_course_time":      newCourseTime,
			"existing_course_time": existingCourseTime,
		},
	}
}

// NewCourseOfferingNotFoundError creates an error for missing course offerings
func NewCourseOfferingNotFoundError(courseOfferingID string) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrCourseOfferingNotFound,
		Message: "Course offering not found",
		Details: map[string]interface{}{
			"course_offering_id": courseOfferingID,
		},
	}
}

// NewInvalidCourseDataError creates an error for invalid course offering data
func NewInvalidCourseDataError(field, reason string) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrInvalidCourseData,
		Message: fmt.Sprintf("Invalid course offering: %s %s", field, reason),
		Details: map[string]interface{}{
			"invalid_field": field,
			"reason":        reason,
		},
	}
}

// NewInvalidTimestampError creates an error for invalid timestamps
func NewInvalidTimestampError(context string) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrInvalidTimestamp,
		Message: fmt.Sprintf("Invalid timestamp: %s", context),
		Details: map[string]interface{}{
			"context": context,
		},
	}
}

// NewDatabaseOperationError creates an error for database operation failures
func NewDatabaseOperationError(operation string, cause error) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrDatabaseOperation,
		Message: fmt.Sprintf("Database operation failed: %s", operation),
		Details: map[string]interface{}{
			"operation":    operation,
			"cause_error":  cause.Error(),
		},
	}
}

// NewTransactionFailedError creates an error for transaction failures
func NewTransactionFailedError(cause error) *EnrollmentError {
	return &EnrollmentError{
		Type:    ErrTransactionFailed,
		Message: "Transaction failed during enrollment process",
		Details: map[string]interface{}{
			"cause_error": cause.Error(),
		},
	}
}

// IsEnrollmentError checks if an error is an EnrollmentError
func IsEnrollmentError(err error) bool {
	_, ok := err.(*EnrollmentError)
	return ok
}

// GetEnrollmentErrorType extracts the error type from an EnrollmentError
func GetEnrollmentErrorType(err error) (EnrollmentErrorType, bool) {
	if enrollmentErr, ok := err.(*EnrollmentError); ok {
		return enrollmentErr.Type, true
	}
	return "", false
}

// IsBusinessRuleViolation checks if the error is a business rule violation
func IsBusinessRuleViolation(err error) bool {
	if enrollmentErr, ok := err.(*EnrollmentError); ok {
		switch enrollmentErr.Type {
		case ErrDuplicateEnrollment, ErrCapacityExceeded, ErrScheduleConflict:
			return true
		}
	}
	return false
}

// IsDataValidationError checks if the error is a data validation error
func IsDataValidationError(err error) bool {
	if enrollmentErr, ok := err.(*EnrollmentError); ok {
		switch enrollmentErr.Type {
		case ErrCourseOfferingNotFound, ErrInvalidCourseData, ErrInvalidTimestamp:
			return true
		}
	}
	return false
}

// IsSystemError checks if the error is a system-level error
func IsSystemError(err error) bool {
	if enrollmentErr, ok := err.(*EnrollmentError); ok {
		switch enrollmentErr.Type {
		case ErrDatabaseOperation, ErrTransactionFailed:
			return true
		}
	}
	return false
}