package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Message types
const (
	MessageTypeTaskUpdate    = "task_update"
	MessageTypeTaskCreate    = "task_create"
	MessageTypeTaskDelete    = "task_delete"
	MessageTypeNotebookUpdate = "notebook_update"
	MessageTypeNotebookCreate = "notebook_create"
	MessageTypeNotebookDelete = "notebook_delete"
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	UserID  uint        `json:"user_id,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	Conn   *websocket.Conn
	UserID uint
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. UserID: %d. Total clients: %d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("Client disconnected. UserID: %d. Total clients: %d", client.UserID, len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			messageData, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				h.mu.RUnlock()
				continue
			}

			for client := range h.clients {
				// Only send to clients with the same UserID or if UserID is 0 (broadcast to all)
				if message.UserID == 0 || client.UserID == message.UserID {
					select {
					case client.Send <- messageData:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

// BroadcastToUser sends a message to all clients of a specific user
func (h *Hub) BroadcastToUser(userID uint, messageType string, data interface{}) {
	message := &Message{
		Type:   messageType,
		Data:   data,
		UserID: userID,
	}
	h.Broadcast(message)
}
