package inventory

import "github.com/google/uuid"

type ItemInstance struct {
	ID         uuid.UUID
	TemplateID int
	Quantity   int // only for stackable items
	Durability int
}
