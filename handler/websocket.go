package handler 

import (
	"net/http"
	"strconv"

	"darkoo/models"
	hub "darkoo/websocket"

	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
)

func WebSocketHandler(c *gin.Context) {
	id := c.Param("group_id")
	groupId, _ := strconv.Atoi(id)

	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)

	if err != nil {
		c.JSON(http.StatusInternalServerError,  gin.H{"errro": "Failed to upgrade connection"})
		return
	}

	client := &hub.Client{Conn: conn, GroupId: uint(groupId)}
	hub.Hub.Register <- client

	defer func() {
		hub.Hub.Unregister <- client
	}()

	for {
		var message models.Message

		if err := conn.ReadJSON(&message); err != nil {
			break
		}
		hub.Hub.Broadcast <- &message
	}
}