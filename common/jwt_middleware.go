package common

import (
	"net/http"
	"siakad-poc/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, BaseResponse[any]{
					Status: StatusError,
					Error: &BaseResponseError{
						Message:   "Authorization header required",
						Details:   []string{"missing Authorization header"},
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      c.Request().RequestURI,
					},
				})
			}

			// Extract token from "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, BaseResponse[any]{
					Status: StatusError,
					Error: &BaseResponseError{
						Message:   "Invalid authorization header format",
						Details:   []string{"authorization header must be 'Bearer <token>'"},
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      c.Request().RequestURI,
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
				return c.JSON(http.StatusUnauthorized, BaseResponse[any]{
					Status: StatusError,
					Error: &BaseResponseError{
						Message:   "Invalid token",
						Details:   []string{err.Error()},
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      c.Request().RequestURI,
					},
				})
			}

			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, BaseResponse[any]{
					Status: StatusError,
					Error: &BaseResponseError{
						Message:   "Token is not valid",
						Details:   []string{"token validation failed"},
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      c.Request().RequestURI,
					},
				})
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, BaseResponse[any]{
					Status: StatusError,
					Error: &BaseResponseError{
						Message:   "Invalid token claims",
						Details:   []string{"unable to parse token claims"},
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      c.Request().RequestURI,
					},
				})
			}

			// Add user information to context
			c.Set(StudentIDKey, claims.UserID)
			c.Set(UserRoleKey, claims.Role)

			return next(c)
		}
	}
}