package inventory

import (
	"errors"
	"sync"
)

var (
	ErrInventoryFull    = errors.New("inventory is full")
	ErrTemplateNotFound = errors.New("template not found")
)

type Inventory struct {
	mu       sync.RWMutex
	capacity int
	items    []*ItemInstance
}

func NewInventory(capacity int) *Inventory {
	return &Inventory{
		capacity: capacity,
		items:    make([]*ItemInstance, 0, capacity),
	}
}

// this mf needs more work
func (inv *Inventory) Add(item *ItemInstance) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	template, templateExists := GetTemplate(item.TemplateID)
	if !templateExists {
		return ErrTemplateNotFound
	}
	lastItem := inv.findByTemplateID(item.TemplateID)
	if lastItem != nil {
		if template.Stackable && (lastItem.Quantity > template.MaxStack) {
			emptyQuantity := template.MaxStack - lastItem.Quantity
			if emptyQuantity > item.Quantity {
				lastItem.Quantity = lastItem.Quantity + item.Quantity
			}
		}
	}
	if len(inv.items) >= inv.capacity {
		return ErrInventoryFull
	}

	inv.items = append(inv.items, item)
	return nil
}

func (inv *Inventory) findByTemplateID(templateID int) *ItemInstance {
	for _, item := range inv.items {
		if item.TemplateID == templateID {
			return item
		}
	}
	return nil
}
