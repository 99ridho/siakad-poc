package modules

import "github.com/gofiber/fiber/v2"

type RoutableModule interface {
	SetupRoutes(fiber *fiber.App, prefix string)
}
