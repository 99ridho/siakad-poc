package handlers

import (
	"siakad-poc/common"
	"siakad-poc/modules/auth/usecases"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type LoginHandler struct {
	usecase *usecases.LoginUseCase
}

type LoginRequestData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type LoginResponseData struct {
	AccessToken string `json:"access_token"`
}

func NewLoginHandler(usecase *usecases.LoginUseCase) *LoginHandler {
	return &LoginHandler{usecase: usecase}
}

func (h *LoginHandler) HandleLogin(c *fiber.Ctx) error {
	requestID := c.Get(fiber.HeaderXRequestID)
	clientIP := c.IP()

	var loginRequest LoginRequestData
	err := c.BodyParser(&loginRequest)
	if err != nil {
		log.Error().
			Stack().
			Err(err).
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.OriginalURL()).
			Str("method", c.Method()).
			Msg("Failed to parse login request body")

		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse login request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	// Validate request data
	if validationErrors := common.ValidateStruct(&loginRequest); validationErrors != nil {
		log.Warn().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", loginRequest.Email).
			Strs("validation_errors", validationErrors).
			Str("path", c.OriginalURL()).
			Msg("Login validation failed")

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

	token, err := h.usecase.Login(c.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Error().
			Stack().
			Err(err).
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", loginRequest.Email).
			Str("path", c.OriginalURL()).
			Msg("Login failed")

		return c.Status(fiber.StatusInternalServerError).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot proceed login",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.BaseResponse[LoginResponseData]{
		Status: common.StatusSuccess,
		Data: &LoginResponseData{
			AccessToken: token,
		},
	})
}
