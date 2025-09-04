package handlers

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/modules/auth/usecases"
	"time"

	"github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	usecase *usecases.RegisterUseCase
}

type RegisterRequestData struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RegisterResponseData struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func NewRegisterHandler(usecase *usecases.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{usecase: usecase}
}

func (h *RegisterHandler) HandleRegister(c echo.Context) error {
	var registerRequest RegisterRequestData
	err := c.Bind(&registerRequest)
	if err != nil {
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

	// Validate password confirmation
	if registerRequest.Password != registerRequest.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Password confirmation does not match",
				Details:   []string{"password and confirm_password must be the same"},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.Request().RequestURI,
			},
		})
	}

	userID, err := h.usecase.Register(c.Request().Context(), registerRequest.Email, registerRequest.Password)
	if err != nil {
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