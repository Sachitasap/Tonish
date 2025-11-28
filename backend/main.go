package main

import (
	"log"
	"tonish/backend/ai"
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
	
	// Initialize AI client
	log.Println("Initializing AI client...")
	ai.Initialize()
	
	// Check AI health
	if err := ai.Client.HealthCheck(); err != nil {
		log.Printf("Warning: AI service not available: %v", err)
		log.Println("The app will run but AI features will be disabled")
	} else {
		log.Printf("AI service connected: %s (model: %s)", ai.Client.BaseURL, ai.Client.Model)
	}

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
	log.Println("Server starting on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
