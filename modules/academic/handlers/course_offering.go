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
		log.Error().Err(err).Msg("Failed to get course offerings")
		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Internal server error",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
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
	var req usecases.CreateCourseOfferingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Invalid request payload",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	if err := common.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   err,
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	response, err := h.useCase.CreateCourseOffering(c.Request().Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create course offering")
		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to create course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	return c.JSON(http.StatusOK, common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleUpdateCourseOffering(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Course offering ID is required",
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	var req usecases.UpdateCourseOfferingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Invalid request payload",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	if err := common.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Validation failed",
				Details:   err,
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	response, err := h.useCase.UpdateCourseOffering(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "course offering not found" {
			return c.JSON(http.StatusNotFound, common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
					Path:      c.Request().URL.Path,
				},
			})
		}

		log.Error().Err(err).Msg("Failed to update course offering")
		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to update course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	return c.JSON(http.StatusOK, common.BaseResponse[usecases.CourseOfferingIDResponse]{
		Status: common.StatusSuccess,
		Data:   &response,
	})
}

func (h *CourseOfferingHandler) HandleDeleteCourseOffering(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Course offering ID is required",
				Details:   []string{"ID parameter is missing"},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	err := h.useCase.DeleteCourseOffering(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "course offering not found" {
			return c.JSON(http.StatusNotFound, common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Course offering not found",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
					Path:      c.Request().URL.Path,
				},
			})
		}

		log.Error().Err(err).Msg("Failed to delete course offering")
		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Failed to delete course offering",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      c.Request().URL.Path,
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}