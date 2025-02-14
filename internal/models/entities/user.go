package entities

// User представляет информацию о пользователе.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Balance  int    `json:"coins"`
}
