package middlewares

import (
	"siakad-poc/common"
	"siakad-poc/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   int64  `json:"role"`
	jwt.RegisteredClaims
}

const (
	StudentIDKey = "student_id"
	UserRoleKey  = "user_role"
)

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Authorization header required",
					Details:   []string{"missing Authorization header"},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Invalid authorization header format",
					Details:   []string{"authorization header must be 'Bearer <token>'"},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKeyType
			}
			return []byte(config.CurrentConfig.JWT.Secret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Invalid token",
					Details:   []string{err.Error()},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Token is not valid",
					Details:   []string{"token validation failed"},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(common.BaseResponse[any]{
				Status: common.StatusError,
				Error: &common.BaseResponseError{
					Message:   "Invalid token claims",
					Details:   []string{"unable to parse token claims"},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Path:      c.OriginalURL(),
				},
			})
		}

		// Add user information to context
		c.Locals(StudentIDKey, claims.UserID)
		c.Locals(UserRoleKey, claims.Role)

		return c.Next()
	}
}
