package websocket

import (
	"net/http"
	"log"

	"darkoo/middleware"
	"strconv"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity; adjust as needed for production
	},
}

// HandleWebSocket handles incoming WebSocket requests
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	userDetails := r.Context().Value("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Extract user and group IDs from query params (or headers, depending on your setup)
	userID := userDetails.(*middleware.User).ID
	groupID := r.URL.Query().Get("groupId")

	if groupID == "" {
		http.Error(w, "userId and groupId are required", http.StatusBadRequest)
		return
	}

	// Create a new client and register it with the Hub
	userId := strconv.Itoa(int(userID))
	client := &Client{
		ID:     userId,
		Group:  groupID,
		Socket: conn,
		Send:   make(chan []byte),
	}
	hub.Register <- client

	// Start read and write pumps for the client
	go client.ReadPump(hub)
	go client.WritePump()
}
