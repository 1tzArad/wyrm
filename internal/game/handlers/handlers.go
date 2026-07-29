package handlers

import (
	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/network"
)

func RegisterHandlers(registry *network.Registery, world *game.World) {
	registry.Register("chat", chatHandler(world))
	registry.Register("move", moveHandler(world))
}
