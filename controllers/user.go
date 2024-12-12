package controllers

import (
	"github.com/chajaykrishna/go-ecommerce/database"
	"github.com/chajaykrishna/go-ecommerce/models"
	"github.com/chajaykrishna/go-ecommerce/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func FetchUserDetails(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(jwt.MapClaims)

	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}
	userId := uint(userIdFloat)

	var user models.User
	// fetch user from db
	if err := database.DB.Where("id", userId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "error while fetching user from db",
				"err":     err,
			},
		})
	}

	var userResponse types.UserResponse
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Phone = user.Phone
	userResponse.Username = user.Username

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": userResponse,
		},
	})
}
