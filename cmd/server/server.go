package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1tzArad/wyrm/internal/auth"
	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/game/handlers"
	"github.com/1tzArad/wyrm/internal/network"
	"github.com/1tzArad/wyrm/internal/storage"
	postgres_repository "github.com/1tzArad/wyrm/internal/storage/postgres/repository"
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

	dbcfg := storage.CreateConfig("postgres", "")
	db, err := storage.Open(*dbcfg)
	if err != nil {
		log.Fatal("Failed to open database connection!", "err", err.Error())
		return
	}

	r := gin.New()

	hub := network.NewHub()
	go hub.Run()

	registery := network.NewRegistery()

	playerRepo := postgres_repository.NewPlayerRepository(db)

	world := game.CreateWorld(hub, playerRepo)

	handlers.RegisterHandlers(registery, world)

	ctx := context.Background()

	go world.RunLoop()
	go world.RunAutoSave(ctx)

	hub.OnPlayerDisconnect = func(playerid uuid.UUID) {
		p, ok := world.GetPlayer(playerid)
		if !ok {
			return
		}
		if err := playerRepo.SavePlayer(context.Background(), p); err != nil {
			log.Errorf("failed to save player %s on disconnect: %v", playerid, err)
		}
		world.RemovePlayer(playerid)
	}

	// registering middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ws", WSHandler(hub, world, registery))

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Warn("shutting down, saving all players...")
		world.SaveAllPlayers(context.Background())
		os.Exit(0)
	}()

	r.Run(":8080")

}

func WSHandler(hub *network.Hub, world *game.World, registery *network.Registery) gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate provided token
		tokenString := c.Query("token")
		user_id, err := auth.ValidateJWT(tokenString)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "INVALID_TOKEN", "invalid or expired token!")
			return
		}

		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade connection!", "err", err.Error())
			response.InternalFail(c)
			return
		}

		player, err := world.LoadOrCreatePlayerForUser(c.Request.Context(), user_id)

		if err != nil {
			log.Error("failed to load/create player", "err", err)
			conn.Close()
			response.InternalFail(c)
			return
		}

		client := network.NewClient(conn, hub, player.ID, registery)

		hub.Register(client)

		go client.ReadPump()
		go client.WritePump()
	}
}
