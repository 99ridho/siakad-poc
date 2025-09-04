package middlewares

import (
	"siakad-poc/common"
	"siakad-poc/constants"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ShouldBeAccessedByRoles(expectedRoles []constants.RoleType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals(UserRoleKey).(constants.RoleType)
		if slices.Contains(expectedRoles, role) {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
			Status: common.StatusError,
			Error: &common.BaseResponseError{
				Message:   "Invalid role",
				Details:   []string{},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Path:      c.OriginalURL(),
			},
		})
	}
}
