package websocket

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"darkoo/models"

	"github.com/gorilla/websocket"
)

// Message represents the structure of incoming WebSocket messages
type Message struct {
	Action        string `json:"action"`         // Action type (joinGroup or sendMessage)
	ContentType   string `json:"contentType"`    // Message content type
	AttachmentURL string `json:"attachmentUrl"`  // Attachment URL (if any)
	GroupID       string `json:"groupId"`        // Group ID
	UserID        string `json:"userId"`         // User ID
	Content       string `json:"content"`        // Message content
}

// Client represents a WebSocket client connection
type Client struct {
	ID     string          // User ID
	Username string
	Group  string          // Group ID the client is part of
	Socket *websocket.Conn // WebSocket connection
	Send   chan []byte     // Channel to send messages to the client
}

// ReadPump listens for incoming messages from the WebSocket connection
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		switch msg.Action {
		case "joinGroup":
			// Handle joining a group
			userId, _ := strconv.Atoi(msg.UserID)
			groupId, _ := strconv.Atoi(msg.GroupID)
			if err := hub.UserService.JoinGroup(userId, groupId); err != nil {
				log.Printf("Failed to join group: %v", err)
				continue
			}
			c.Group = msg.GroupID
			log.Printf("User %s joined group %s", msg.UserID, msg.GroupID)

			joinNotification := map[string]string{
				"action":  "groupNotification",
				"message": "User " + msg.UserID + " has joined group " + msg.GroupID,
			}

			notificationJSON, err := json.Marshal(joinNotification)
			if err != nil {
				log.Printf("Failed to marshal join notification: %v", err)
				continue
			}
			hub.Broadcast <- notificationJSON

		case "sendMessage":
			// Handle sending a message
			groupId, _ := strconv.Atoi(msg.GroupID)
			userId, _ := strconv.Atoi(msg.UserID)
			var attachmentURL *string

			if msg.AttachmentURL != "" {
    			attachmentURL = &msg.AttachmentURL
			}
			


			newMessage := models.Message{
				ContentType:   msg.ContentType,
				AttachmentUrl: attachmentURL,
				GroupId:       uint(groupId),
				UserId:        uint(userId),
				Content:       msg.Content,
			}

			sentMessage, err := hub.MessageService.SendMessage(&newMessage);
			
			if err != nil {
				log.Printf("Failed to save message: %v", err)
				continue
			}
			broadcastMessage, err := json.Marshal(sentMessage)

    		if err != nil {
        		log.Printf("Failed to marshal message: %v", err)
        		continue
    		}

    // Broadcast the saved message to all clients
    hub.Broadcast <- broadcastMessage

		default:
			log.Printf("Unknown action: %s", msg.Action)
		}
	}
}

// WritePump sends messages to the client via WebSocket
func (c *Client) WritePump() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}
	}
}




// func (c *Client) ReadPump(hub *Hub) {
// 	defer func() {
// 		hub.Unregister <- c
// 		c.Socket.Close()
// 	}()

// 	for {
// 		// Read the incoming message
// 		_, message, err := c.Socket.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			break
// 		}

// 		// Assuming your message format is in JSON and includes groupID, userID, contentType, etc.
// 		var msg Message
// 		if err := json.Unmarshal(message, &msg); err != nil {
// 			log.Println("Error unmarshalling message:", err)
// 			continue
// 		}

// 		// If it's a join group action
// 		if msg.Action == "join" {
// 			// Join group via UserService
// 			if err := hub.UserService.JoinGroup(msg.UserID, msg.GroupID); err != nil {
// 				log.Println("Error joining group:", err)
// 				continue
// 			}
// 		} else if msg.Action == "sendMessage" {
// 			// Save message via MessageService
// 			if err := hub.MessageService.SaveMessage(msg); err != nil {
// 				log.Println("Error saving message:", err)
// 				continue
// 			}
// 		}

// 		// Broadcast the message to all clients in the group
// 		hub.Broadcast <- message
// 	}
// }

// func (c *Client) WritePump() {
// 	// Writes messages to the WebSocket
// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			if !ok {
// 				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			c.Socket.WriteMessage(websocket.TextMessage, message)
// 		}
// 	}
// }







// package main

// import (
// 	"log"
// 	"net/http"
// 	"yourapp/message"  // Replace with the actual package path
// 	"yourapp/user"     // Replace with the actual package path
// 	"yourapp/websocket" // Replace with the actual package path
// )

// func main() {
// 	// Initialize your services
// 	messageService := message.NewService()  // Initialize the MessageService
// 	userService := user.NewService()        // Initialize the UserService

// 	// Create the WebSocket Hub
// 	hub := websocket.NewHub(messageService, userService)

// 	// Start the Hub in a goroutine
// 	go hub.Start()

// 	// Set up the WebSocket handler
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		websocket.WebSocketHandler(hub, w, r)
// 	})

// 	// Start the HTTP server
// 	log.Println("Server started on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
