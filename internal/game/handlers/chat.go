package handlers

import (
	"encoding/json"

	"github.com/1tzArad/wyrm/internal/game"
	"github.com/1tzArad/wyrm/internal/network"
	"github.com/charmbracelet/log"
)

func ChatHandler(world *game.World) network.HandlerFunc {
	return func(c *network.Context) {
		log.Debug("received new chat")
		var data game.ChatPayload
		if err := c.Bind(&data); err != nil {
			return
		}

		chatType := data.Type
		switch chatType {
		case game.GLOBAL:
			handleGlobalChat(&data, c, world)
		default:
			log.Error("invalid chat type!", "type", chatType)
		}
	}
}

func handleGlobalChat(data *game.ChatPayload, c *network.Context, world *game.World) {
	log.Debug("received new global chat")
	outMsg := network.Message{
		Type: "chat",
	}

	payload := game.ChatBroadcastMessage{
		PlayerUUID: c.PlayerUUID.String(),
		Text:       data.Text,
	}

	playloadBytes, _ := json.Marshal(payload)
	outMsg.Payload = playloadBytes

	finalBytes, _ := json.Marshal(outMsg)

	world.Hub.Broadcast(finalBytes)
}
