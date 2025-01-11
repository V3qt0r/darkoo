package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"darkoo/models"
)

type WebSocketHub struct {
	Groups 		map[uint]map[*websocket.Conn]bool
	Register 	chan *Client
	Unregister 	chan *Client
	Broadcast 	chan *models.Message
	mu sync.Mutex
}

type Client struct {
	Conn 	*websocket.Conn
	GroupId uint
}

var Hub = &WebSocketHub{Groups: make(map[uint]map[*websocket.Conn]bool),
	Register: 	make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast: 	make(chan *models.Message),
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.Groups[client.GroupId]; !ok {
				h.Groups[client.GroupId] = make(map[*websocket.Conn]bool)
			}
		h.Groups[client.GroupId][client.Conn] = true
		h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Groups[client.GroupId]; ok {
				if _, ok := h.Groups[client.GroupId][client.Conn]; ok {
						delete(h.Groups[client.GroupId], client.Conn)
						client.Conn.Close()
				}
			if len(h.Groups[client.GroupId]) == 0 {
				delete(h.Groups, client.GroupId)
			}
		} 

		h.mu.Unlock()
		
		case message := <-h.Broadcast: 
		h.mu.Lock()
			if clients, ok := h.Groups[message.GroupId]; ok {
				for conn := range clients {
					if err := conn.WriteJSON(message); err != nil {
						conn.Close()
						delete(clients, conn)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}
