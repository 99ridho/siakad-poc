package main

import (
	"context"
	"siakad-poc/config"
	"siakad-poc/db/repositories"
	"siakad-poc/modules/auth/handlers"
	"siakad-poc/modules/auth/usecases"

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

	usersRepository := repositories.NewDefaultUserRepository(pool)
	loginUseCase := usecases.NewLoginUseCase(usersRepository)
	loginHandler := handlers.NewLoginHandler(loginUseCase)

	e := echo.New()
	e.POST("/login", loginHandler.HandleLogin)

	e.Logger.Fatal(e.Start(":8880"))
}
