package handlers

import (
	"encoding/json"
	"errors"

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
			c.ReplyError("invalid_body", "request body is invald!")
			return
		}

		chatType := data.Type
		switch chatType {
		case game.GLOBAL:
			handleGlobalChat(&data, c, world)
		case game.PRIVATE:
			handlePrivateChat(&data, c, world)
		default:
			c.ReplyError("invalid_chat", "chat type is not valid!")
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
		c.ReplyError("empty_target", "target should not be empty!")
		log.Error("no target")
		return
	}
	target, err := uuid.Parse(targetStr)
	if err != nil {
		c.ReplyError("invalid_target", "target is not valid!")
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

	if err := world.Broadcaster.SendToPlayer(target, finalBytes); err != nil {
		if errors.Is(err, network.ErrPlayerNotConnected) {
			c.ReplyError("player_offline", "target player is offline!")
			return
		}
		if errors.Is(err, network.ErrPlayerChannelFull) {
			c.ReplyError("target_channel_full", "target player's channel is full!")
			return
		}
		c.ReplyError("internal_error", err.Error())
		return
	}
}
