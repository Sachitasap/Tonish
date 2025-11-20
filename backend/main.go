package main

import (
	"log"
	"tonish/backend/database"
	"tonish/backend/middleware"
	"tonish/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	database.Connect()
	database.Migrate()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Tonish API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(middleware.CORS())

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Tonish API is running",
			"version": "1.0",
		})
	})

	// Setup routes
	routes.Setup(app)

	// Start server
	log.Println("Server starting on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
