package handlers

import (
	"encoding/json"

	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/network"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func chatHandler(world *game.World) network.HandlerFunc {
	return func(c *network.Context) {
		log.Debug("received new chat")
		var data game.ChatPayload
		if err := c.Bind(&data); err != nil {
			// server response
			return
		}

		chatType := data.Type
		switch chatType {
		case game.GLOBAL:
			handleGlobalChat(&data, c, world)
		case game.PRIVATE:
			handlePrivateChat(&data, c, world)
		default:
			// server response
			log.Error("invalid chat type!", "type", chatType)
		}
	}
}

func handleGlobalChat(data *game.ChatPayload, c *network.Context, world *game.World) {
	outMsg := network.Message{
		Type: "chat",
	}

	payload := game.ChatReceivePayload{
		PlayerUUID: c.PlayerUUID.String(),
		Text:       data.Text,
	}

	playloadBytes, _ := json.Marshal(payload)
	outMsg.Payload = playloadBytes

	finalBytes, _ := json.Marshal(outMsg)

	world.Broadcaster.Broadcast(finalBytes)
}

func handlePrivateChat(data *game.ChatPayload, c *network.Context, world *game.World) {
	targetStr := data.Target
	if targetStr == "" {
		// server response
		log.Error("no target")
		return
	}
	target, err := uuid.Parse(targetStr)
	if err != nil {
		// server response
		log.Error("invalid target")
		return
	}

	outMsg := network.Message{
		Type: "chat",
	}

	payload := game.ChatReceivePayload{
		PlayerUUID: c.PlayerUUID.String(),
		Text:       data.Text,
	}

	payloadBytes, _ := json.Marshal(payload)
	outMsg.Payload = payloadBytes

	finalBytes, _ := json.Marshal(outMsg)

	world.Broadcaster.SendToPlayer(target, finalBytes)
}
