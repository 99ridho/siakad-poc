package middlewares

import (
	"net/http"
	"siakad-poc/common"
	"siakad-poc/constants"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
)

func ShouldBeAccessedByRoles(expectedRoles []constants.RoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get(UserRoleKey).(constants.RoleType)
			if slices.Contains(expectedRoles, role) {
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Invalid role",
					Details:   []string{},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.Request().RequestURI,
				},
			})
		}
	}
}
