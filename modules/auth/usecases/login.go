package usecases

import (
	"context"
	"siakad-poc/config"
	"siakad-poc/constants"
	"siakad-poc/db/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	repository repositories.UserRepository
}

type JWTClaims struct {
	UserID string             `json:"user_id"`
	Role   constants.RoleType `json:"role"`
	jwt.RegisteredClaims
}

func NewLoginUseCase(repository repositories.UserRepository) *LoginUseCase {
	return &LoginUseCase{repository: repository}
}

func (u *LoginUseCase) Login(ctx context.Context, email, password string) (string, error) {
	// Get user by email
	user, err := u.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		return "", errors.Wrap(err, "failed to get user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Convert role from pgtype.Numeric to RoleType (int64)
	var userRole constants.RoleType
	err = user.Role.Scan(&userRole)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse user role")
	}

	// Generate JWT token
	claims := JWTClaims{
		UserID: user.ID.String(),
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.CurrentConfig.JWT.Secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return tokenString, nil
}
