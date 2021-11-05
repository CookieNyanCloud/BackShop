package domain

type User struct {
	ID       string `json:"-" db:"id"`
	Email    string `form:"email" json:"email" binding:"required" db:"email"`
	Password string `form:"password" json:"password" binding:"required" db:"password_hash"`
}

type Verification struct {
	ID    string `json:"id" binding:"required" db:"id"`
	Code  string `db:"codes"`
	State bool   `json:"state"`
}
