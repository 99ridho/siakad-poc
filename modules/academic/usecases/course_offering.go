package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"siakad-poc/common"
	"siakad-poc/db/repositories"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseOfferingResponse struct {
	ID          string    `json:"id"`
	CourseName  string    `json:"course_name"`
	CourseCode  string    `json:"course_code"`
	SectionCode string    `json:"section_code"`
	Capacity    int32     `json:"capacity"`
	StartTime   time.Time `json:"start_time"`
}

type CreateCourseOfferingRequest struct {
	CourseID    string    `json:"course_id" validate:"required"`
	SemesterID  string    `json:"semester_id" validate:"required"`
	SectionCode string    `json:"section_code" validate:"required"`
	Capacity    int32     `json:"capacity" validate:"required,min=1"`
	StartTime   time.Time `json:"start_time" validate:"required"`
}

type UpdateCourseOfferingRequest struct {
	CourseID    string    `json:"course_id" validate:"required"`
	SemesterID  string    `json:"semester_id" validate:"required"`
	SectionCode string    `json:"section_code" validate:"required"`
	Capacity    int32     `json:"capacity" validate:"required,min=1"`
	StartTime   time.Time `json:"start_time" validate:"required"`
}

type CourseOfferingIDResponse struct {
	ID string `json:"id"`
}

type CourseOfferingUseCase struct {
	repo repositories.AcademicRepository
}

func NewCourseOfferingUseCase(repo repositories.AcademicRepository) *CourseOfferingUseCase {
	return &CourseOfferingUseCase{
		repo: repo,
	}
}

func (uc *CourseOfferingUseCase) GetCourseOfferingsWithPagination(ctx context.Context, page, pageSize int) ([]CourseOfferingResponse, *common.PaginationMetadata, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	courseOfferings, err := uc.repo.GetCourseOfferingsWithPagination(ctx, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}

	totalRecords, err := uc.repo.CountCourseOfferings(ctx)
	if err != nil {
		return nil, nil, err
	}

	var responses []CourseOfferingResponse
	for _, co := range courseOfferings {
		startTime := time.Time{}
		if co.CourseOfferingStartTime.Valid {
			startTime = co.CourseOfferingStartTime.Time
		}

		responses = append(responses, CourseOfferingResponse{
			ID:          uuidToString(co.CourseOfferingID),
			CourseName:  co.CourseName,
			CourseCode:  co.CourseCode,
			SectionCode: co.SectionCode,
			Capacity:    co.Capacity,
			StartTime:   startTime,
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	pagination := &common.PaginationMetadata{
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: int(totalRecords),
		TotalPages:   totalPages,
	}

	return responses, pagination, nil
}

func (uc *CourseOfferingUseCase) CreateCourseOffering(ctx context.Context, req CreateCourseOfferingRequest) (CourseOfferingIDResponse, error) {
	courseOffering, err := uc.repo.CreateCourseOffering(ctx, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime)
	if err != nil {
		return CourseOfferingIDResponse{}, err
	}

	return CourseOfferingIDResponse{
		ID: uuidToString(courseOffering.ID),
	}, nil
}

func (uc *CourseOfferingUseCase) UpdateCourseOffering(ctx context.Context, id string, req UpdateCourseOfferingRequest) (CourseOfferingIDResponse, error) {
	courseOffering, err := uc.repo.UpdateCourseOffering(ctx, id, req.SemesterID, req.CourseID, req.SectionCode, req.Capacity, req.StartTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CourseOfferingIDResponse{}, errors.New("course offering not found")
		}
		return CourseOfferingIDResponse{}, err
	}

	return CourseOfferingIDResponse{
		ID: uuidToString(courseOffering.ID),
	}, nil
}

func (uc *CourseOfferingUseCase) DeleteCourseOffering(ctx context.Context, id string) error {
	_, err := uc.repo.DeleteCourseOffering(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("course offering not found")
		}
		return err
	}
	return nil
}

func uuidToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid.Bytes[0:4],
		uuid.Bytes[4:6],
		uuid.Bytes[6:8],
		uuid.Bytes[8:10],
		uuid.Bytes[10:16])
}
