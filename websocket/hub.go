package websocket

import (
	"log"
	"darkoo/models"

	"encoding/json"
)

// Hub manages active WebSocket clients and broadcasts messages
type Hub struct {
	Clients    map[string]*Client      // A map of client IDs to clients
	Register   chan *Client            // Channel to register new clients
	Unregister chan *Client           // Channel to unregister clients
	Broadcast  chan []byte            // Channel for broadcasting messages to all clients
	MessageService models.IMessageService
	UserService	   models.IUserService             // Interface to interact with the database
}

// NewHub creates a new instance of Hub
func NewHub(messageService models.IMessageService, userService models.IUserService) *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		MessageService: messageService,  // Pass the concrete implementation
		UserService:    userService,     // Pass the concrete implementation
	}
}

// Start runs the Hub's main loop, listening for Register, Unregister, and Broadcast events
func (h *Hub) Start() {
    for {
        select {
        case client := <-h.Register:
            // Register a new client
            h.Clients[client.ID] = client
            log.Printf("New client registered: %s", client.ID)

        case client := <-h.Unregister:
            // Unregister a client
            delete(h.Clients, client.ID)
            log.Printf("Client unregistered: %s", client.ID)

        case message := <-h.Broadcast:
            // Broadcast the message to all clients in the same group
            var msg Message
            if err := json.Unmarshal(message, &msg); err != nil {
                log.Printf("Invalid broadcast message: %v", err)
                continue
            }

            for _, client := range h.Clients {
                if client.Group == msg.GroupID { // Send only to clients in the same group
                    select {
                    case client.Send <- message:
                    default:
                        close(client.Send)
                        delete(h.Clients, client.ID)
                    }
                }
            }
        }
    }
}

