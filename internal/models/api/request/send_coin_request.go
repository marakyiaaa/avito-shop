package request

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
