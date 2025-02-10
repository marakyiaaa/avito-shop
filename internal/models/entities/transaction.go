package entities

import "time"

type Transaction struct {
	ID         int       `json:"id"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Amount     int       `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
