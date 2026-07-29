package game

import (
	"encoding/json"
	"time"

	"github.com/1tzArad/wyrm/internal/network"
)

const TickRate = time.Second * 10 // it must be 60 tick per second

func (w *World) RunLoop() {
	ticker := time.NewTicker(TickRate)
	defer ticker.Stop()

	for range ticker.C {
		w.tick()
	}
}

func (w *World) tick() {
	snapshot := w.Snapshot()
	w.broadcastSnapshot(snapshot)
}

func (w *World) broadcastSnapshot(snapshot []PlayerState) {
	msg := network.Message{Type: "snapshot"}

	payloadBytes, err := json.Marshal(snapshot)
	if err != nil {
		return
	}
	msg.Payload = payloadBytes

	finalBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	w.Broadcaster.Broadcast(finalBytes)
}
