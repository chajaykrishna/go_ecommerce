package middlewares

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chajaykrishna/go-ecommerce/config"
	"github.com/chajaykrishna/go-ecommerce/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

func GenerateJwtToken(user *models.User) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		return "", errors.New("jwt_secret is not set")
	}
	now := time.Now()
	accessTokenClaims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   strconv.Itoa(int(user.ID)),
			Issuer:    "go-commerce",
			Audience:  []string{"test"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(config.Jwt_Expiry)),
		},
		UserID: user.ID,
		Email:  user.Email,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "nil", err
	}

	return accessTokenString, nil
}

func ValidateJwtToken(c *fiber.Ctx) error {

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "JWT secret not configured",
			},
		})
	}

	header := c.Get("Authorization")
	if header == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	parts := strings.Split(header, " ")

	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "Missing Authorization header",
			},
		})
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwt_secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Store token claims in context
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}
	c.Locals("userClaims", claims)
	return c.Next()
}
