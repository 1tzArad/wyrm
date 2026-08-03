package inventory

var templates = map[int]ItemTemplate{
	1: {ID: 1, Name: "Iron Sword", Description: "A Sword Made by a retard", Type: ItemTypeWeapon, Stackable: false, MaxDurability: 100},
	2: {ID: 2, Name: "Wooden Sword", Description: "A Sword made out of wood", Type: ItemTypeWeapon, Stackable: false, MaxDurability: 50},
	3: {ID: 3, Name: "Wooden Axe", Description: "A Axe made out of wood", Type: ItemTypeWeapon, Stackable: false, MaxDurability: 50},
}

func GetTemplate(id int) (ItemTemplate, bool) {
	t, ok := templates[id]
	return t, ok
}
