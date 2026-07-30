package game

import (
	"context"
	"time"

	"github.com/1tzArad/wyrm/internal/player"
	"github.com/charmbracelet/log"
)

const autoSaveInterval = 30 * time.Second

func (w *World) RunAutoSave(ctx context.Context) {
	ticker := time.NewTicker(autoSaveInterval)

	for {
		select {
		case <-ticker.C:
			w.SaveAllPlayers(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *World) SaveAllPlayers(ctx context.Context) {
	w.mu.Lock()
	playersCopy := make([]*player.Player, 0, len(w.players))
	for _, p := range w.players {
		playersCopy = append(playersCopy, p)
	}
	w.mu.Unlock()

	for _, p := range playersCopy {
		if err := w.store.SavePlayer(ctx, p); err != nil {
			log.Errorf("failed to save player %s: %v", p.ID, err)
		}
	}
}
