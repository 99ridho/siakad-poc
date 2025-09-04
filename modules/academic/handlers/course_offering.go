package handlers

import (
	"siakad-poc/common"
	"siakad-poc/modules/academic/usecases"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type CourseOfferingHandler struct {
	useCase *usecases.CourseOfferingUseCase
}

func NewCourseOfferingHandler(useCase *usecases.CourseOfferingUseCase) *CourseOfferingHandler {
	return &CourseOfferingHandler{
		useCase: useCase,
	}
}

func (h *CourseOfferingHandler) HandleListCourseOfferings(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

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

	courseOfferings, pagination, err := h.useCase.GetCourseOfferingsWithPagination(c.Context(), page, pageSize)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Int("page", page).
			Int("page_size", pageSize).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Failed to get course offerings")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Internal server error",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.PaginatedBaseResponse[[]usecases.CourseOfferingResponse]{
		BaseResponse: common.BaseResponse[[]usecases.CourseOfferingResponse]{
			Status: common.StatusSuccess,
			Data:   &courseOfferings,
		},
		Paging: pagination,
	})
}

func (h *CourseOfferingHandler) HandleCreateCourseOffering(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	var req usecases.CreateCourseOfferingRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Failed to parse create course offering request body")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
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
			Str("path", c.OriginalURL()).
			Msg("Create course offering validation failed")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   validationErrors,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	response, err := h.useCase.CreateCourseOffering(c.Context(), req)
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
			Str("path", c.OriginalURL()).
			Msg("Failed to create course offering")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to create course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleUpdateCourseOffering(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	id := c.Params("id")
	if id == "" {
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
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	var req usecases.UpdateCourseOfferingRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Failed to parse update course offering request body")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
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
			Str("path", c.OriginalURL()).
			Msg("Update course offering validation failed")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   validationErrors,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	response, err := h.useCase.UpdateCourseOffering(c.Context(), id, req)
	if err != nil {
		if err.Error() == "course offering not found" {
			log.Warn().
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("course_offering_id", id).
				Str("path", c.OriginalURL()).
				Msg("Course offering not found for update")

			return c.Status(fiber.StatusNotFound).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
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
			Str("path", c.OriginalURL()).
			Msg("Failed to update course offering")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to update course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleDeleteCourseOffering(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	id := c.Params("id")
	if id == "" {
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
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	err := h.useCase.DeleteCourseOffering(c.Context(), id)
	if err != nil {
		if err.Error() == "course offering not found" {
			log.Warn().
				Str("request_id", requestID).
				Str("client_ip", clientIP).
				Str("course_offering_id", id).
				Str("path", c.OriginalURL()).
				Msg("Course offering not found for deletion")

			return c.Status(fiber.StatusNotFound).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("course_offering_id", id).
			Str("path", c.OriginalURL()).
			Msg("Failed to delete course offering")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to delete course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
