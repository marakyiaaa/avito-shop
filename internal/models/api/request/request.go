package request

// AuthRequest тело запроса для аутентификации.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=2,max=50"`
}

// SendCoinRequest тело запроса для отправки монет.
type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
