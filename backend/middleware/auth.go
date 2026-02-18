package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return []byte(secret)
}

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}
	
	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}
	
	tokenString := parts[1]
	
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return getJWTSecret(), nil
	})
	
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}
	
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}
	
	// Store user ID in context
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	c.Locals("email", claims["email"].(string))
	
	return c.Next()
}
