package main

import (
	"flag"
	"os"
	"time"

	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/network"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

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

	world := game.CreateWorld(hub)

	// registering middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ws", network.Handler(hub, world))

	r.Run(":8080")
}
