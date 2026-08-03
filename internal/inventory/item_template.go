package inventory

type ItemType string

const (
	ItemTypeWeapon ItemType = "weapons"
)

type ItemTemplate struct {
	ID            int
	Name          string
	Description   string
	Type          ItemType
	Stackable     bool
	MaxStack      int
	MaxDurability int
}
