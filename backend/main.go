package main

import (
	"log"
	"os"
	"tonish/backend/database"
	"tonish/backend/middleware"
	"tonish/backend/routes"
	ws "tonish/backend/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	database.Connect()
	database.Migrate()
	database.NormalizeTaskTypes()
	database.SeedDefaultUser()

	// Initialize WebSocket hub
	ws.Initialize()

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
			"message": "Tonish API is running - REGISTRATION DISABLED",
			"version": "1.0.1",
		})
	})

	// Setup routes
	routes.Setup(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
