package entities

// Inventory представляет инвентарь пользователя.
type Inventory struct {
	ID       int
	UserID   int
	ItemType string
	Quantity int
}
