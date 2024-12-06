package controllers

import (
	"errors"
	"strings"

	"github.com/chajaykrishna/go-ecommerce/database"
	"github.com/chajaykrishna/go-ecommerce/middlewares"
	"github.com/chajaykrishna/go-ecommerce/models"
	"github.com/chajaykrishna/go-ecommerce/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username" validate:"required,min=3,alphanum"`
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"omitempty,min=3"`
}

type UserResponse struct {
	Username string `json:"username" validate:"required,min=3,alphanum"`
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
}

var validate = validator.New()

func Signup(c *fiber.Ctx) error {

	var request User
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
	// normalize data
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Name = strings.TrimSpace(request.Name)
	request.Username = strings.ToLower(strings.TrimSpace(request.Username))

	// check if the username, phone, email already exist
	var existingUser User
	if err := database.DB.Where("email=? OR username=? OR phone=?",
		request.Email, request.Username, request.Phone).First(&existingUser).Error; err == nil {
		return sendErrorResponse(c, fiber.StatusConflict, errors.New("User with username/mail/phone already exist"))
	}

	// Hash password before storing in DB
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, errors.New("error processing the registration"))
	}
	request.Password = hashedPassword

	if err := database.DB.Create(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "signup failed",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user signup successful",
	})
}

func ValidateUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	var userCount int64
	if err := database.DB.Model(&models.User{}).Where("username=?", username).Count(&userCount).Error; err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, errors.New("internal error"))
	}
	if userCount > 0 {
		return sendErrorResponse(c, fiber.StatusConflict, errors.New("username already taken"))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Username available",
		},
	})
}

func Login(c *fiber.Ctx) error {
	var userRequest types.UserSigninRquest
	if err := c.BodyParser(&userRequest); err != nil {
		return sendErrorResponse(c, fiber.StatusBadGateway, errors.New("invalid request body format"))
	}

	// validate the user body
	if err := validate.Struct(&userRequest); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("invalid request data"))
	}

	var user models.User
	if err := database.DB.Where("email=?", userRequest.Email).First(&user).Error; err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("invalid credentials"))
	}

	// validate user given password with actual password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, errors.New("invalid credentials"))
	}

	userResponse := types.UserResponse{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	accessToken, err := middlewares.GenerateJwtToken(&user)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, errors.New("internal error"))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "user login successful",
			"user":    userResponse,
			"token":   accessToken,
		},
	})

}

// internal functions
func sendErrorResponse(c *fiber.Ctx, statuscode int, err error) error {
	return c.Status(statuscode).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"message": "error",
			"error":   err.Error(),
		},
	})
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err

}
