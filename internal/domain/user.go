package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"-" db:"id"`
	Email      string    `form:"email" json:"email" binding:"required" db:"email"`
	Password   string    `form:"password" json:"password" binding:"required" db:"password_hash"`
}

type Verification struct {
	Email string `form:"email" json:"email" binding:"required" db:"email"`
	Code  string `db:"codes"`
	State bool   `json:"state"`
}
