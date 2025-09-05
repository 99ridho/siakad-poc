package auth

import (
	"siakad-poc/db/repositories"
	"siakad-poc/modules/auth/handlers"
	"siakad-poc/modules/auth/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	userRepository  repositories.UserRepository
	loginUseCase    *usecases.LoginUseCase
	loginHandler    *handlers.LoginHandler
	registerUseCase *usecases.RegisterUseCase
	registerHandler *handlers.RegisterHandler
}

func NewModule(pool *pgxpool.Pool) *AuthModule {
	usersRepository := repositories.NewDefaultUserRepository(pool)

	loginUseCase := usecases.NewLoginUseCase(usersRepository)
	loginHandler := handlers.NewLoginHandler(loginUseCase)

	registerUseCase := usecases.NewRegisterUseCase(usersRepository)
	registerHandler := handlers.NewRegisterHandler(registerUseCase)

	return &AuthModule{
		userRepository:  usersRepository,
		loginUseCase:    loginUseCase,
		loginHandler:    loginHandler,
		registerUseCase: registerUseCase,
		registerHandler: registerHandler,
	}
}

func (m *AuthModule) SetupRoutes(fiberApp *fiber.App, prefix string) {
	authRoutes := fiberApp.Group(prefix)
	authRoutes.Post("/login", m.loginHandler.HandleLogin)
	authRoutes.Post("/register", m.registerHandler.HandleRegister)
}
