package game

import (
	"sync"

	"github.com/1tzArad/wyrm/internal/network"
	"github.com/1tzArad/wyrm/internal/player"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Broadcaster interface {
	Broadcast(message []byte)
	SendToPlayer(playerUUID uuid.UUID, message []byte) error
}

type World struct {
	mu sync.RWMutex

	players     map[uuid.UUID]*player.Player
	hub         *network.Hub
	Broadcaster Broadcaster
}

func CreateWorld(hub *network.Hub) *World {
	return &World{
		players:     make(map[uuid.UUID]*player.Player),
		hub:         hub,
		Broadcaster: hub,
	}
}

func (w *World) CreatePlayer(uuid uuid.UUID) *player.Player {
	w.mu.Lock()
	defer w.mu.Unlock()
	p := &player.Player{
		UUID: uuid,
		X:    0,
		Y:    0,
	}

	w.players[uuid] = p
	return p
}

func (w *World) RemovePlayer(uuid uuid.UUID) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.players, uuid)
}

func (w *World) GetPlayer(uuid uuid.UUID) (*player.Player, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	p, ok := w.players[uuid]
	return p, ok
}

func (w *World) MovePlayer(uuid uuid.UUID, direction Direction) {
	w.mu.Lock()
	defer w.mu.Unlock()

	p, ok := w.players[uuid]
	if !ok {
		return
	}

	switch direction {
	case UP:
		p.Y--
	case DOWN:
		p.Y++
	case LEFT:
		p.X--
	case RIGHT:
		p.X++
	}
	log.Debugf("New player move to %s", direction)
	log.Debugf("new x: %f | new y: %f", p.X, p.Y)
}

func (w *World) Snapshot() []PlayerState {
	w.mu.Lock()
	defer w.mu.Unlock()

	states := make([]PlayerState, 0, len(w.players))
	for _, p := range w.players {
		states = append(states, PlayerState{
			UUID:   p.UUID,
			X:      p.X,
			Y:      p.Y,
			Health: p.Health,
		})
	}
	return states
}
