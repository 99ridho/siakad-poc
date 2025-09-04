package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/modules/academic/usecases"
	"time"

	"github.com/labstack/echo/v4"
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
	// Extract course offering ID from URL parameter
	courseOfferingID := c.Param("id")
	if courseOfferingID == "" {
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
	studentIDInterface := c.Get(common.StudentIDKey)
	if studentIDInterface == nil {
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
	err := h.enrollmentUseCase.EnrollStudentToCourseOffering(c.Request().Context(), studentID, courseOfferingID)
	if err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusBadRequest

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
