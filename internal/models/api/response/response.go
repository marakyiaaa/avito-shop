package response

// InventoryItem представляет предмет в инвентаре пользователя.
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinTransaction представляет транзакцию с монетами.
type CoinTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

// CoinHistory представляет историю монет пользователя.
type CoinHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

// InfoResponse тело ответа для получения информации.
type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

// AuthResponse тело ответа для аутентификации.
type AuthResponse struct {
	Token string `json:"token"`
}

// ErrorResponse тело ответа для ошибок.
type ErrorResponse struct {
	Errors string `json:"errors"`
}
