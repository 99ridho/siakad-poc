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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
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

	e := echo.New()

	// Auth routes (unprotected)
	e.POST("/login", loginHandler.HandleLogin)
	e.POST("/register", registerHandler.HandleRegister)

	// Academic routes (protected with JWT middleware)
	academicGroup := e.Group("/academic")
	academicGroup.Use(middlewares.JWT())
	academicGroup.POST(
		"/course-offering/:id/enroll",
		enrollmentHandler.HandleCourseEnrollment,
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
	)

	e.Logger.Fatal(e.Start(":8880"))
}
