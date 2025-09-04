package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/modules/auth/usecases"
	"time"

	"github.com/labstack/echo/v4"
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

func (h *RegisterHandler) HandleRegister(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	var registerRequest RegisterRequestData
	err := c.Bind(&registerRequest)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Failed to parse register request body")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse register request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
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
			Str("path", c.Request().RequestURI).
			Msg("Registration validation failed")

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

	userID, err := h.usecase.Register(c.Request().Context(), registerRequest.Email, registerRequest.Password)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", registerRequest.Email).
			Str("path", c.Request().RequestURI).
			Msg("Registration failed")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot proceed registration",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.JSON(http.StatusCreated, common.BaseResponse[RegisterResponseData]{
		Status: common.StatusSuccess,
		Data: &RegisterResponseData{
			UserID:  userID,
			Message: "User registered successfully",
		},
	})
}
