package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/modules/auth/usecases"
	"time"

	"github.com/labstack/echo/v4"
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

func (h *LoginHandler) HandleLogin(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Request().Header.Get("X-Request-ID")
	}
	clientIP := c.RealIP()

	var loginRequest LoginRequestData
	err := c.Bind(&loginRequest)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("path", c.Request().RequestURI).
			Str("method", c.Request().Method).
			Msg("Failed to parse login request body")

		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot parse login request body",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
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
			Str("path", c.Request().RequestURI).
			Msg("Login validation failed")

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

	token, err := h.usecase.Login(c.Request().Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("email", loginRequest.Email).
			Str("path", c.Request().RequestURI).
			Msg("Login failed")

		return c.JSON(http.StatusInternalServerError, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Cannot proceed login",
				Details:   []string{err.Error()},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	return c.JSON(http.StatusOK, common.BaseResponse[LoginResponseData]{
		Status: common.StatusSuccess,
		Data: &LoginResponseData{
			AccessToken: token,
		},
	})
}
