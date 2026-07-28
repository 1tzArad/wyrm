package game

import (
	"github.com/1tzArad/wyrm/internal/network"
	"github.com/1tzArad/wyrm/internal/player"
	"github.com/google/uuid"
)

type World struct {
	players map[uuid.UUID]*player.Player
	hub     *network.Hub
}



func CreateWorld(hub *network.Hub) *World {
	return &World{
		players: make(map[uuid.UUID]*player.Player),
		hub:     hub,
	}
}

func (w *World) CreatePlayer(uuid uuid.UUID) *player.Player {
	p := &player.Player{
		UUID: uuid,
		X:    0,
		Y:    0,
	}

	w.players[uuid] = p
	return p
}
