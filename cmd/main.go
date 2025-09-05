package main

import (
	"context"
	"siakad-poc/config"
	"siakad-poc/modules"
	"siakad-poc/modules/academic"
	"siakad-poc/modules/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
}

func main() {
	pool, err := pgxpool.New(context.Background(), config.CurrentConfig.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create database pool")
	}

	routePrefixToModuleMapping := map[string]modules.RoutableModule{
		"/auth":     auth.NewModule(pool),
		"/academic": academic.NewModule(pool),
	}

	app := fiber.New()
	app.Use(
		cors.New(),
		helmet.New(),
		recover.New(),
		logger.New(),
		healthcheck.New(healthcheck.Config{
			LivenessEndpoint:  "/live",
			ReadinessEndpoint: "/ready",
		}),
	)

	for pfx, module := range routePrefixToModuleMapping {
		module.SetupRoutes(app, pfx)
	}

	log.Fatal().Err(app.Listen(config.CurrentConfig.App.Addr)).Msg("Failed to start server")
}
