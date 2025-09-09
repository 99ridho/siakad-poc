package handlers

import (
	"siakad-poc/common"
	"siakad-poc/middlewares"
	"siakad-poc/modules/academic/usecases"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type CourseEnrollmentHandler struct {
	enrollmentUseCase *usecases.CourseEnrollmentUseCase
}

type EnrollmentResponseData struct {
	Message           string    `json:"message"`
	StudentID         string    `json:"student_id"`
	CourseOfferingID  string    `json:"course_offering_id"`
	EnrollmentTime    time.Time `json:"enrollment_time"`
	Status            string    `json:"status"`
}

func NewEnrollmentHandler(enrollmentUseCase *usecases.CourseEnrollmentUseCase) *CourseEnrollmentHandler {
	return &CourseEnrollmentHandler{
		enrollmentUseCase: enrollmentUseCase,
	}
}

func (h *CourseEnrollmentHandler) HandleCourseEnrollment(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	// Extract course offering ID from URL parameter
	courseOfferingID := c.Params("id")
	if courseOfferingID == "" {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Course offering ID missing from URL parameter")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Course offering ID is required",
				Details:   []string{"course offering ID must be provided in URL path"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	// Extract student ID from JWT context (set by middleware)
	studentIDInterface := c.Locals(middlewares.StudentIDKey)
	if studentIDInterface == nil {
		log.Error().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", courseOfferingID).
			Str("path", c.OriginalURL()).
			Msg("Student ID not found in JWT token context")

		return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Student ID not found in token",
				Details:   []string{"authentication token does not contain student ID"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	studentID, ok := studentIDInterface.(string)
	if !ok {
		log.Error().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", courseOfferingID).
			Interface("student_id_raw", studentIDInterface).
			Str("path", c.OriginalURL()).
			Msg("Student ID from JWT token is not in valid string format")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Invalid student ID format",
				Details:   []string{"student ID from token is not in valid format"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	// Call use case to enroll student
	err := h.enrollmentUseCase.EnrollStudent(c.Context(), studentID, courseOfferingID)
	if err != nil {
		// Determine appropriate HTTP status code and user-friendly message based on error type
		statusCode := fiber.StatusBadRequest
		userMessage := "Enrollment failed"
		errorDetails := []string{err.Error()}

		// Handle domain-specific errors with better UX
		if enrollmentErr, ok := err.(*usecases.EnrollmentError); ok {
			switch enrollmentErr.Type {
			case usecases.ErrDuplicateEnrollment:
				statusCode = fiber.StatusConflict
				userMessage = "You are already enrolled in this course"
				errorDetails = []string{"Duplicate enrollment detected. You cannot enroll in the same course offering twice."}

			case usecases.ErrCapacityExceeded:
				statusCode = fiber.StatusConflict
				userMessage = "Course is full"
				errorDetails = []string{"This course offering has reached its maximum capacity. Please try a different section or contact the academic office."}

			case usecases.ErrScheduleConflict:
				statusCode = fiber.StatusConflict
				userMessage = "Schedule conflict detected"
				errorDetails = []string{"The selected course conflicts with your existing class schedule. Please choose a different time slot."}

			case usecases.ErrCourseOfferingNotFound:
				statusCode = fiber.StatusNotFound
				userMessage = "Course offering not found"
				errorDetails = []string{"The requested course offering does not exist or may have been cancelled."}

			case usecases.ErrInvalidCourseData:
				statusCode = fiber.StatusBadRequest
				userMessage = "Invalid course information"
				errorDetails = []string{"There is an issue with the course offering data. Please contact the academic office."}

			case usecases.ErrDatabaseOperation:
				statusCode = fiber.StatusInternalServerError
				userMessage = "System temporarily unavailable"
				errorDetails = []string{"A technical issue occurred. Please try again later or contact support if the problem persists."}

			case usecases.ErrTransactionFailed:
				statusCode = fiber.StatusInternalServerError
				userMessage = "Enrollment could not be processed"
				errorDetails = []string{"A system error prevented enrollment completion. Please try again."}

			default:
				// Keep default values for unknown enrollment errors
				userMessage = "Enrollment failed"
				errorDetails = []string{enrollmentErr.Error()}
			}
		}

		// Log the enrollment failure with structured context
		logEvent := log.Error().
			Stack().
			Err(err).
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("student_id", studentID).
			Str("course_offering_id", courseOfferingID).
			Str("path", c.OriginalURL()).
			Int("http_status", statusCode).
			Str("user_message", userMessage)

		// Add enrollment error details if available
		if enrollmentErr, ok := err.(*usecases.EnrollmentError); ok {
			logEvent = logEvent.
				Str("error_type", string(enrollmentErr.Type)).
				Interface("error_details", enrollmentErr.Details).
				Bool("is_business_rule_violation", usecases.IsBusinessRuleViolation(err)).
				Bool("is_system_error", usecases.IsSystemError(err))
		}

		logEvent.Msg("Course enrollment failed")

		return c.Status(statusCode).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   userMessage,
				Details:   errorDetails,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	// Log successful enrollment
	log.Info().
		Str("request_id", requestID).
		Str("client_ip", clientIP).
		Str("student_id", studentID).
		Str("course_offering_id", courseOfferingID).
		Str("path", c.OriginalURL()).
		Msg("Course enrollment successful")

	// Return enhanced success response with enrollment details
	enrollmentTime := time.Now().UTC()
	return c.Status(fiber.StatusCreated).JSON(common.BaseResponse[EnrollmentResponseData]{
		Status: common.StatusSuccess,
		Data: &EnrollmentResponseData{
			Message:          "Successfully enrolled in course offering",
			StudentID:        studentID,
			CourseOfferingID: courseOfferingID,
			EnrollmentTime:   enrollmentTime,
			Status:           "enrolled",
		},
	})
}
