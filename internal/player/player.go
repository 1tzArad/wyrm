package player

import (
	"github.com/1tzArad/wyrm/internal/inventory"
	"github.com/google/uuid"
)

type Player struct {
	ID           uuid.UUID
	X, Y         int32
	Health, Mana int32

	Inventory *inventory.Inventory
}
