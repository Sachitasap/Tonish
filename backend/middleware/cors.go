package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	
	if allowedOrigins == "" {
		// Default to allow all in development (no credentials with wildcard)
		return cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
			AllowCredentials: false,
		})
	}

	// Production: specific origins with credentials
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
