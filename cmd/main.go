package main

import (
	"context"
	"os"
	"os/signal"
	"siakad-poc/config"
	"siakad-poc/modules"
	"siakad-poc/modules/academic"
	"siakad-poc/modules/auth"
	"syscall"
	"time"

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
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection pool
	pool, err := pgxpool.New(ctx, config.CurrentConfig.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create database pool")
	}

	// Mapping HTTP route prefix to relevant module
	routePrefixToModuleMapping := map[string]modules.RoutableModule{
		"/auth":     auth.NewModule(pool),
		"/academic": academic.NewModule(pool),
	}

	// Initialize HTTP handler library
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

	// Setup routes per module
	for pfx, module := range routePrefixToModuleMapping {
		module.SetupRoutes(app, pfx)
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Info().Str("address", config.CurrentConfig.App.Addr).Msg("Starting server")
		if err := app.Listen(config.CurrentConfig.App.Addr); err != nil {
			log.Error().Err(err).Msg("Server failed to start or stopped")
		}
	}()

	log.Info().Msg("Server started successfully. Press Ctrl+C to gracefully shutdown")

	// Block until a signal is received
	<-quit
	log.Info().Msg("Graceful shutdown initiated...")

	// Create a deadline for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Gracefully shutdown the server
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server shutdown gracefully")
	}

	// Close database connection pool
	log.Info().Msg("Closing database connections...")
	pool.Close()
	log.Info().Msg("Database connections closed")

	log.Info().Msg("Application shutdown completed")
}
