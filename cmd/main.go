package main

import (
	"context"
	"siakad-poc/config"
	"siakad-poc/constants"
	"siakad-poc/db/repositories"
	"siakad-poc/middlewares"
	academicHandlers "siakad-poc/modules/academic/handlers"
	academicUsecases "siakad-poc/modules/academic/usecases"
	authHandlers "siakad-poc/modules/auth/handlers"
	authUsecases "siakad-poc/modules/auth/usecases"

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

	// Repositories
	usersRepository := repositories.NewDefaultUserRepository(pool)
	academicRepository := repositories.NewDefaultAcademicRepository(pool)

	// Auth module
	loginUseCase := authUsecases.NewLoginUseCase(usersRepository)
	loginHandler := authHandlers.NewLoginHandler(loginUseCase)

	registerUseCase := authUsecases.NewRegisterUseCase(usersRepository)
	registerHandler := authHandlers.NewRegisterHandler(registerUseCase)

	// Academic module
	enrollmentUseCase := academicUsecases.NewCourseEnrollmentUseCase(academicRepository)
	enrollmentHandler := academicHandlers.NewEnrollmentHandler(enrollmentUseCase)

	courseOfferingUseCase := academicUsecases.NewCourseOfferingUseCase(academicRepository)
	courseOfferingHandler := academicHandlers.NewCourseOfferingHandler(courseOfferingUseCase)

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

	// Auth routes (unprotected)
	app.Post("/login", loginHandler.HandleLogin)
	app.Post("/register", registerHandler.HandleRegister)

	// Academic routes (protected with JWT middleware)
	academicGroup := app.Group("/academic")
	academicGroup.Use(middlewares.JWT())
	academicGroup.Post(
		"/course-offering/:id/enroll",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
		enrollmentHandler.HandleCourseEnrollment,
	)

	// Course offering CRUD routes (Admin and Koorprodi only)
	academicGroup.Get(
		"/course-offering",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		courseOfferingHandler.HandleListCourseOfferings,
	)
	academicGroup.Post(
		"/course-offering",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		courseOfferingHandler.HandleCreateCourseOffering,
	)
	academicGroup.Put(
		"/course-offering/:id",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		courseOfferingHandler.HandleUpdateCourseOffering,
	)
	academicGroup.Delete(
		"/course-offering/:id",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		courseOfferingHandler.HandleDeleteCourseOffering,
	)

	log.Fatal().Err(app.Listen(config.CurrentConfig.App.Addr)).Msg("Failed to start server")
}
