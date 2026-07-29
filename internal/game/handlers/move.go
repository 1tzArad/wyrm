package handlers

import (
	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/network"
)

func moveHandler(world *game.World) network.HandlerFunc {
	return func(c *network.Context) {
		var data game.MovePayload
		if err := c.Bind(&data); err != nil {
			c.ReplyError("invalid_body", "invalid body!")
			return
		}

		if !data.Direction.IsValid() {
			c.ReplyError("invalid_move", "invalid move!")
			return
		}

		world.MovePlayer(c.PlayerUUID, data.Direction)
	}
}
