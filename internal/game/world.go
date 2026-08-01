package game

import (
	"context"
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
	store       PlayerStore
}

func CreateWorld(hub *network.Hub, store PlayerStore) *World {
	return &World{
		players:     make(map[uuid.UUID]*player.Player),
		hub:         hub,
		Broadcaster: hub,
		store:       store,
	}
}

func (w *World) CreatePlayer(uuid uuid.UUID) *player.Player {
	w.mu.Lock()
	defer w.mu.Unlock()
	p := &player.Player{
		ID: uuid,
		X:  0,
		Y:  0,
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
	log.Debugf("new x: %d | new y: %d", p.X, p.Y)
}

func (w *World) Snapshot() []PlayerState {
	w.mu.Lock()
	defer w.mu.Unlock()

	states := make([]PlayerState, 0, len(w.players))
	for _, p := range w.players {
		states = append(states, PlayerState{
			UUID:   p.ID,
			X:      p.X,
			Y:      p.Y,
			Health: p.Health,
		})
	}
	return states
}

func (w *World) LoadOrCreatePlayerForUser(ctx context.Context, userID uuid.UUID) (*player.Player, error) {
	existingPlayer, err := w.store.GetByUserID(ctx, userID)
	if err == nil {
		w.players[existingPlayer.ID] = existingPlayer
		return existingPlayer, nil
	}

	newPlayer, err := w.store.Create(ctx, userID)
	if err != nil {
		return nil, err
	}

	w.players[newPlayer.ID] = newPlayer
	return newPlayer, nil
}
