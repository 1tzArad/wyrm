package network

import (
	"net/http"

	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/pkg/response"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func Handler(hub *Hub, world *game.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade connection!", "err", err.Error())
			response.InternalFail(c)
			return
		}

		playerUUID := uuid.New()
		world.CreatePlayer(playerUUID)

		client := NewClient(conn, hub, playerUUID)

		go client.ReadPump()
		go client.WritePump()
	}
}
