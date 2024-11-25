package controllers

import (
	"github.com/chajaykrishna/go-ecommerce/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,alphanum"`
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

var validate = validator.New()

func Signup(c *fiber.Ctx) error {

	var request SignupRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
			"errors":  err.Error(),
		})
	}

	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
			"errors":  err.Error(),
		})
	}

	if err := database.DB.Create(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "signup failed",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user signup successful",
		"user":    request,
	})
}
