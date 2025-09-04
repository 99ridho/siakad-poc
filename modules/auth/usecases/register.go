package usecases

import (
	"context"
	"regexp"
	"siakad-poc/db/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
	repository repositories.UserRepository
}

const (
	DefaultStudentRole = 3
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewRegisterUseCase(repository repositories.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{repository: repository}
}

func (u *RegisterUseCase) Register(ctx context.Context, email, password string) (string, error) {
	// Validate email format
	if !emailRegex.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	// Check if user already exists
	_, err := u.repository.GetUserByEmail(ctx, email)
	if err == nil {
		return "", errors.New("email already registered")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "failed to check existing user")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	// Create user with default student role
	user, err := u.repository.CreateUser(ctx, email, string(hashedPassword), DefaultStudentRole)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	return user.ID.String(), nil
}
