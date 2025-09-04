package handlers

import (
	"siakad-poc/common"
	"siakad-poc/modules/auth/usecases"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type RegisterHandler struct {
	usecase *usecases.RegisterUseCase
}

type RegisterRequestData struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type RegisterResponseData struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func NewRegisterHandler(usecase *usecases.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{usecase: usecase}
}

func (h *RegisterHandler) HandleRegister(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	var registerRequest RegisterRequestData
	err := c.BodyParser(&registerRequest)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Failed to parse register request body")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse register request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	// Validate request data
	if validationErrors := common.ValidateStruct(&registerRequest); validationErrors != nil {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", registerRequest.Email).
			Strs("validation_errors", validationErrors).
			Str("path", c.OriginalURL()).
			Msg("Registration validation failed")

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

	userID, err := h.usecase.Register(c.Context(), registerRequest.Email, registerRequest.Password)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", registerRequest.Email).
			Str("path", c.OriginalURL()).
			Msg("Registration failed")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot proceed registration",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.BaseResponse[RegisterResponseData]{
		Status: common.StatusSuccess,
		Data: &RegisterResponseData{
			UserID:  userID,
			Message: "User registered successfully",
		},
	})
}
