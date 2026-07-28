package player

import "github.com/google/uuid"

type Player struct {
	UUID         uuid.UUID
	X, Y         float64
	Health, Mana int
}

