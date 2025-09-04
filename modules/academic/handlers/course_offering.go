package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/modules/academic/usecases"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CourseOfferingHandler struct {
	useCase usecases.CourseOfferingUseCase
}

func NewCourseOfferingHandler(useCase usecases.CourseOfferingUseCase) *CourseOfferingHandler {
	return &CourseOfferingHandler{
		useCase: useCase,
	}
}

func (h *CourseOfferingHandler) HandleListCourseOfferings(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	courseOfferings, pagination, err := h.useCase.GetCourseOfferingsWithPagination(c.Request().Context(), page, pageSize)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Int("page", page).
			Int("page_size", pageSize).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Failed to get course offerings")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Internal server error",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.JSON(http.StatusOK, common.PaginatedBaseResponse[[]usecases.CourseOfferingResponse]{
		BaseResponse: common.BaseResponse[[]usecases.CourseOfferingResponse]{
			Status: common.StatusSuccess,
			Data:   &courseOfferings,
		},
		Paging: pagination,
	})
}

func (h *CourseOfferingHandler) HandleCreateCourseOffering(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	var req usecases.CreateCourseOfferingRequest
	if err := c.Bind(&req); err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Failed to parse create course offering request body")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	if validationErrors := common.ValidateStruct(req); validationErrors != nil {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_id", req.CourseID).
			Str("semester_id", req.SemesterID).
			Str("section_code", req.SectionCode).
			Int32("capacity", req.Capacity).
			Strs("validation_errors", validationErrors).
			Str("path", c.Request().RequestURI).
			Msg("Create course offering validation failed")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   validationErrors,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	response, err := h.useCase.CreateCourseOffering(c.Request().Context(), req)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_id", req.CourseID).
			Str("semester_id", req.SemesterID).
			Str("section_code", req.SectionCode).
			Int32("capacity", req.Capacity).
			Str("path", c.Request().RequestURI).
			Msg("Failed to create course offering")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to create course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.JSON(http.StatusOK, common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleUpdateCourseOffering(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	id := c.Param("id")
	if id == "" {
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
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	var req usecases.UpdateCourseOfferingRequest
	if err := c.Bind(&req); err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Failed to parse update course offering request body")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	if validationErrors := common.ValidateStruct(req); validationErrors != nil {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("course_id", req.CourseID).
			Str("semester_id", req.SemesterID).
			Str("section_code", req.SectionCode).
			Int32("capacity", req.Capacity).
			Strs("validation_errors", validationErrors).
			Str("path", c.Request().RequestURI).
			Msg("Update course offering validation failed")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   validationErrors,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	response, err := h.useCase.UpdateCourseOffering(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "course offering not found" {
			log.Warn().
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("course_offering_id", id).
				Str("path", c.Request().RequestURI).
				Msg("Course offering not found for update")

			return c.JSON(http.StatusNotFound, common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.Request().RequestURI,
				},
			})
		}

		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("course_id", req.CourseID).
			Str("semester_id", req.SemesterID).
			Str("section_code", req.SectionCode).
			Int32("capacity", req.Capacity).
			Str("path", c.Request().RequestURI).
			Msg("Failed to update course offering")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to update course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.JSON(http.StatusOK, common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleDeleteCourseOffering(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	id := c.Param("id")
	if id == "" {
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
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	err := h.useCase.DeleteCourseOffering(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "course offering not found" {
			log.Warn().
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("course_offering_id", id).
				Str("path", c.Request().RequestURI).
				Msg("Course offering not found for deletion")

			return c.JSON(http.StatusNotFound, common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.Request().RequestURI,
				},
			})
		}

		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("path", c.Request().RequestURI).
			Msg("Failed to delete course offering")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to delete course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}