package routes

import (
	"github.com/chajaykrishna/go-ecommerce/controllers"
	"github.com/chajaykrishna/go-ecommerce/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {

	app.Post("api/v1/signup", controllers.Signup)

	app.Get("api/v1/validateUsername/:username", controllers.ValidateUsername)

	app.Get("api/v1/login", controllers.Login)

	app.Get("api/v1/myDetails", middlewares.ValidateJwtToken, controllers.FetchUserDetails)

	app.Post("api/v1/addProduct", middlewares.ValidateJwtToken, controllers.AddProduct)

	app.Get("api/v1/products", middlewares.ValidateJwtToken, controllers.GetAllProducts)

	app.Get("api/v1/product/:{id}", middlewares.ValidateJwtToken, controllers.GetProductById)

	app.Put("api/v1/product/:id", middlewares.ValidateJwtToken, controllers.UpdateProduct)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Route not found",
			"message": "The requested route doesn't exist.",
		})

	})

}
