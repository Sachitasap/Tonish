package websocket

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var GlobalHub *Hub

// Initialize sets up the WebSocket hub
func Initialize() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket hub initialized")
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *websocket.Conn) {
	// Get user ID from query params
	userIDParam := c.Query("user_id")
	var userID uint
	if userIDParam != "" {
		// Convert string to uint
		var id int
		if _, err := fmt.Sscanf(userIDParam, "%d", &id); err == nil {
			userID = uint(id)
		}
	}

	client := &Client{
		Conn:   c,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}

	GlobalHub.Register(client)

	// Start goroutine to send messages
	go func() {
		for {
			message, ok := <-client.Send
			if !ok {
				// Channel closed
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		}
	}()

	// Read messages (keep connection alive)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}

	GlobalHub.Unregister(client)
}

// Middleware upgrades HTTP to WebSocket
func Middleware() fiber.Handler {
	return websocket.New(HandleWebSocket)
}
