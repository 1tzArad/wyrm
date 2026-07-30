package player

import "github.com/google/uuid"

type Player struct {
	ID           uuid.UUID
	X, Y         int32
	Health, Mana int32
}
