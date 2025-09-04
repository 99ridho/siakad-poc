package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/middlewares"
	"siakad-poc/modules/academic/usecases"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CourseEnrollmentHandler struct {
	enrollmentUseCase *usecases.CourseEnrollmentUseCase
}

type EnrollmentResponseData struct {
	Message string `json:"message"`
}

func NewEnrollmentHandler(enrollmentUseCase *usecases.CourseEnrollmentUseCase) *CourseEnrollmentHandler {
	return &CourseEnrollmentHandler{
		enrollmentUseCase: enrollmentUseCase,
	}
}

func (h *CourseEnrollmentHandler) HandleCourseEnrollment(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	// Extract course offering ID from URL parameter
	courseOfferingID := c.Param("id")
	if courseOfferingID == "" {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Course offering ID missing from URL parameter")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Course offering ID is required",
				Details:   []string{"course offering ID must be provided in URL path"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	// Extract student ID from JWT context (set by middleware)
	studentIDInterface := c.Get(middlewares.StudentIDKey)
	if studentIDInterface == nil {
		log.Error().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", courseOfferingID).
			Str("path", c.Request().RequestURI).
			Msg("Student ID not found in JWT token context")

		return c.JSON(http.StatusUnauthorized, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Student ID not found in token",
				Details:   []string{"authentication token does not contain student ID"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
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
			Str("path", c.Request().RequestURI).
			Msg("Student ID from JWT token is not in valid string format")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Invalid student ID format",
				Details:   []string{"student ID from token is not in valid format"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	// Call use case to enroll student
	err := h.enrollmentUseCase.EnrollStudent(c.Request().Context(), studentID, courseOfferingID)
	if err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusBadRequest

		// Log the enrollment failure with context
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("student_id", studentID).
			Str("course_offering_id", courseOfferingID).
			Str("path", c.Request().RequestURI).
			Msg("Course enrollment failed")

		// You could implement more sophisticated error type checking here
		// For now, treating all business logic errors as bad request

		return c.JSON(statusCode, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Enrollment failed",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	// Return success response
	return c.JSON(http.StatusCreated, common.BaseResponse[EnrollmentResponseData]{
		Status: common.StatusSuccess,
		Data: &EnrollmentResponseData{
			Message: "Successfully enrolled in course offering",
		},
	})
}
