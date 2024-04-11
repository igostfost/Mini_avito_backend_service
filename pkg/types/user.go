package types

type User struct {
	Id       int    `json:"-" db:"user_id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
}
