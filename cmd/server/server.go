package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/game/handlers"
	"github.com/1tzArad/wyrm/internal/network"
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

func init() {
	var isDebuggingMode bool
	flag.BoolVar(&isDebuggingMode, "debug", false, "debugging mode")
	flag.Parse()

	newLog := log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:      time.DateTime,
		Prefix:          "WYRM",
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	if isDebuggingMode {
		newLog.SetLevel(log.DebugLevel)
	}

	log.SetDefault(newLog)

}

func main() {
	r := gin.New()

	hub := network.NewHub()
	go hub.Run()

	registery := network.NewRegistery()

	world := game.CreateWorld(hub)

	handlers.RegisterHandlers(registery, world)

	// registering middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ws", WSHandler(hub, world, registery))

	r.Run(":8080")
}

func WSHandler(hub *network.Hub, world *game.World, registery *network.Registery) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade connection!", "err", err.Error())
			response.InternalFail(c)
			return
		}

		playerUUID := uuid.New()
		world.CreatePlayer(playerUUID)

		client := network.NewClient(conn, hub, playerUUID, registery)

		hub.Register(client)

		go client.ReadPump()
		go client.WritePump()
	}
}
